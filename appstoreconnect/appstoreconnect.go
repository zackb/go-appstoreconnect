package appstoreconnect

import (
	"compress/gzip"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"gopkg.in/yaml.v2"
)

const (
	baseURL = "https://api.appstoreconnect.apple.com/v1/"
)

var (
	// ErrAuthKeyNotPem the key provided is not a pem file
	ErrAuthKeyNotPem = errors.New("token: AuthKey must be a valid .p8 PEM file")

	// ErrAuthKeyNotECDSA the key provided is not in the expected format (apple)
	ErrAuthKeyNotECDSA = errors.New("token: AuthKey must be of type ecdsa.PrivateKey")

	// ErrNoData the endpoint returned a 404, which means do data in that time range
	ErrNoData = errors.New("no data for date range")
)

// Credentials holder for all the information needed to communicate with appstore connect api
type Credentials struct {
	KeyID        string `yaml:"key_id"`
	IssuerID     string `yaml:"issuer_id"`
	PrivKey      string `yaml:"private_key"`
	VendorNumber string `yaml:"vendor_number"`
}

// Client to use to communicate with the connect api
type Client struct {
	jwt           string
	vendorNumber  string
	client        *http.Client
	SalesReport   *SalesReport
	FinanceReport *FinanceReport

	// TODO token expiry
}

type service struct {
	Path   string
	Params map[string]string
	client *Client
}

// NewClient creates a new connect api client using the provided credential and config information
func NewClient(creds *Credentials) (*Client, error) {
	payload := jwt.StandardClaims{
		Audience:  "appstoreconnect-v1",
		Issuer:    creds.IssuerID,
		ExpiresAt: time.Now().Unix() + 600,
	}

	token := jwt.Token{
		Header: map[string]interface{}{
			"typ": "JWT",
			"alg": "ES256",
			"kid": creds.KeyID,
		},
		Claims: payload,
		Method: jwt.SigningMethodES256,
	}

	key, err := parseP8PrivKey([]byte(creds.PrivKey))
	if err != nil {
		log.Fatal("parse p8 priv key fail")
		return nil, err
	}
	client := new(Client)
	secretStr, err := token.SignedString(key)
	client.jwt = secretStr
	client.vendorNumber = creds.VendorNumber

	// Reuse a single struct instead of allocating one for each service on the heap
	svc := &service{
		client: client,
	}

	client.SalesReport = (*SalesReport)(svc)
	client.FinanceReport = (*FinanceReport)(svc)

	client.initClient()
	return client, err
}

// NewCredentials creates a representation of the credentials needed to communicate with the connect api
// keyId is found in the Users and Access section of the UI: https://appstoreconnect.apple.com/access/api
// issuerId is found in the Users and Access section of the UI: https://appstoreconnect.apple.com/access/api
// privateKey is downloaded from the same screen https://appstoreconnect.apple.com/access/api
func NewCredentials(keyId string, issuerId string, privateKey string) *Credentials {
	return &Credentials{
		KeyID:    keyId,
		IssuerID: issuerId,
		PrivKey:  privateKey,
	}
}

func NewCredentialsFromFile(path string) (*Credentials, error) {
	c := new(Credentials)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	e := yaml.UnmarshalStrict(b, c)
	return c, e
}

func NewClientFromCredentialsFile(path string) (*Client, error) {
	creds, err := NewCredentialsFromFile(path)
	if err != nil {
		return nil, err
	}
	return NewClient(creds)
}

func (c *Client) initClient() {
	t := &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 30 * time.Second,
	}

	c.client = &http.Client{
		Transport: t,
	}
}

func makeUrl(path string) string {
	return baseURL + path
}

// Get make a request to the Apple App Store Connect API
func (c *Client) get(path string, params map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", makeUrl(path), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/a-gzip")
	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Add("Authorization", "Bearer "+c.jwt)

	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}

	q.Add("filter[vendorNumber]", c.vendorNumber)

	req.URL.RawQuery = q.Encode()

	fmt.Println(req.URL)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode > 299 && resp.StatusCode != 404 {
		body, _ := ioutil.ReadAll(resp.Body)
		return body, errors.New(string(body))
	}

	if resp.StatusCode == 404 {
		return nil, ErrNoData
	}

	z, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, err
	}
	defer z.Close()

	return ioutil.ReadAll(z)
}

func parseP8PrivKey(bytes []byte) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode(bytes)
	if block == nil {
		return nil, ErrAuthKeyNotPem
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	switch pk := key.(type) {
	case *ecdsa.PrivateKey:
		return pk, nil
	default:
		return nil, ErrAuthKeyNotECDSA
	}
}
