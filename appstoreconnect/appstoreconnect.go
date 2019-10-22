package appstoreconnect

import (
	"bytes"
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
	"github.com/zackb/go-appstoreconnect/tsv"
	"gopkg.in/yaml.v2"
)

const (
	BaseUrl = "https://api.appstoreconnect.apple.com/v1/"
)

var (
	ErrAuthKeyNotPem   = errors.New("token: AuthKey must be a valid .p8 PEM file")
	ErrAuthKeyNotECDSA = errors.New("token: AuthKey must be of type ecdsa.PrivateKey")
)

type Credentials struct {
	KeyID        string `yaml:"key_id"`
	IssuerID     string `yaml:"issuer_id"`
	PrivKey      string `yaml:"private_key"`
	VendorNumber string `yaml:"vendor_number"`
}

type Client struct {
	jwtToken     string
	vendorNumber string
	client       *http.Client

	// TODO token expiry
}

type service struct {
	Path   string
	Params map[string]string
}

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
	client.jwtToken = secretStr
	client.vendorNumber = creds.VendorNumber
	client.initClient()
	return client, err
}

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
	return BaseUrl + path
}

func (c *Client) GetSalesReport(date time.Time, frequency Frequency, reportType ReportType, reportSubType ReportSubType) ([]*SalesReportResponse, error) {
	b, err := c.Get(NewSalesReport(date, frequency, reportType, reportSubType))
	if err != nil {
		return nil, err
	}

	data := SalesReportResponse{}
	p, err := tsv.NewParser(bytes.NewReader(b), &data)
	if err != nil {
		return nil, err
	}

	ret := []*SalesReportResponse{}
	for {
		eof, err := p.Next()
		if eof {
			break
		}
		if err != nil {
			return nil, err
		}
		ret = append(ret, data.Clone())
	}
	return ret, nil
}

func (c *Client) Get(s *service) ([]byte, error) {
	req, err := http.NewRequest("GET", makeUrl(s.Path), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/a-gzip")
	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Add("Authorization", "Bearer "+c.jwtToken)

	q := req.URL.Query()
	for k, v := range s.Params {
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
	if resp.StatusCode > 299 {
		body, _ := ioutil.ReadAll(resp.Body)
		return body, errors.New(string(body))
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
