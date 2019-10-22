package appstoreconnect

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/zackb/go-appstoreconnect/encoding"
)

// Sales and Trends reports
// https://developer.apple.com/documentation/appstoreconnectapi/download_sales_and_trends_reports

type ReportType string

type ReportSubType string

type SalesReportResponse struct {
	Reports []*SalesReport
}
type SalesReport struct {
	Provider              string `tsv:"Provider"`
	ProviderCountry       string `tsv:"Provider Country"`
	SKU                   string `tsv:"SKU"`
	Developer             string `tsv:"Developer"`
	Title                 string `tsv:"Title"`
	Version               string `tsv:"Version"`
	ProductTypeIdentifier string `tsv:"Product Type Identifier"`
	Units                 int    `tsv:"Units"`
	DeveloperProceeds     string `tsv:"Developer Proceeds"`
	BeginDate             string `tsv:"Begin Date"`
	EndDate               string `tsv:"End Date"`
	CustomerCurrency      string `tsv:"Customer Currency"`
	CountryCode           string `tsv:"Country Code"`
	CurrencyOfProceeds    string `tsv:"Currency of Proceeds"`
	AppleIdentifier       string `tsv:"Apple Identifier"`
	CustomerPrice         string `tsv:"Customer Price"`
	PromoCode             string `tsv:"Promo Code"`
	ParentIdentifier      string `tsv:"Parent Identifier"`
	Subscription          string `tsv:"Subscription"`
	Period                string `tsv:"Period"`
	Category              string `tsv:"Category"`
	CMB                   string `tsv:"CMB"`
	Device                string `tsv:"Device"`
	SupportedPlatforms    string `tsv:"Supported Platforms"`
	ProceedsReason        string `tsv:"Proceeds Reason"`
	PreservedPricing      string `tsv:"Preserved Pricing"`
	Client                string `tsv:"Client"`
	OrderType             string `tsv:"Order Type"`
}

const (
	Daily   Frequency = "DAILY"
	Weekly  Frequency = "WEEKLY"
	Monthly Frequency = "MONTHLY"
	Yearly  Frequency = "YEARLY"
)

const (
	ReportSales ReportType = "SALES"
)

const (
	SubReportSummary ReportSubType = "SUMMARY"
)

func NewSalesReport(date time.Time, frequency Frequency, reportType ReportType, reportSubType ReportSubType) *service {
	return &service{
		Path: "salesReports",
		Params: map[string]string{
			"filter[frequency]":     frequency.String(),
			"filter[reportDate]":    timeToReportDate(date, frequency),
			"filter[reportType]":    reportType.String(),
			"filter[reportSubType]": reportSubType.String(),
		},
	}
}

func (f *Frequency) String() string {
	return string(*f)
}

func (r *ReportType) String() string {
	return string(*r)
}

func (r *ReportSubType) String() string {
	return string(*r)
}

func timeToReportDate(t time.Time, f Frequency) string {
	var format string
	switch f {
	case Daily:
		format = "2006-01-02"
	case Weekly:
		t = t.AddDate(0, 0, -int(t.Weekday()))
		format = "2006-01-02"
	case Monthly:
		format = "2006-01"
	case Yearly:
		format = "2006"
	default:
	}

	return t.Format(format)
}

// Clone deep copy
func (s *SalesReport) Clone() *SalesReport {
	c := SalesReport{}
	c.Provider = s.Provider
	c.ProviderCountry = s.ProviderCountry
	c.SKU = s.SKU
	c.Developer = s.Developer
	c.Title = s.Title
	c.Version = s.Version
	c.ProductTypeIdentifier = s.ProductTypeIdentifier
	c.Units = s.Units
	c.DeveloperProceeds = s.DeveloperProceeds
	c.BeginDate = s.BeginDate
	c.EndDate = s.EndDate
	c.CustomerCurrency = s.CustomerCurrency
	c.CountryCode = s.CountryCode
	c.CurrencyOfProceeds = s.CurrencyOfProceeds
	c.AppleIdentifier = s.AppleIdentifier
	c.CustomerPrice = s.CustomerPrice
	c.PromoCode = s.PromoCode
	c.ParentIdentifier = s.ParentIdentifier
	c.Subscription = s.Subscription
	c.Period = s.Period
	c.Category = s.Category
	c.CMB = s.CMB
	c.Device = s.Device
	c.SupportedPlatforms = s.SupportedPlatforms
	c.ProceedsReason = s.ProceedsReason
	c.PreservedPricing = s.PreservedPricing
	c.Client = s.Client
	c.OrderType = s.OrderType
	return &c
}

func (s *SalesReport) GetHeader() []string {
	t := reflect.TypeOf(SalesReport{})
	tags := []string{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("tsv")
		tags = append(tags, tag)
	}

	return tags
}

func (s *SalesReport) Values() []string {
	vals := []string{}
	t := reflect.ValueOf(s).Elem()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		vals = append(vals, fmt.Sprintf("%v", f.Interface()))
	}

	return vals
}

func (s *SalesReportResponse) ToJson() ([]byte, error) {
	return json.Marshal(s)
}

func (s *SalesReportResponse) ToCsv() ([]byte, error) {
	b := bytes.Buffer{}
	buf := bufio.NewWriter(&b)

	w := csv.NewWriter(buf)
	d := SalesReport{}
	headers := d.GetHeader()
	if err := w.Write(headers); err != nil {
		return nil, err
	}

	for _, r := range s.Reports {
		values := r.Values()
		if err := w.Write(values); err != nil {
			return nil, err
		}
	}
	buf.Flush()
	return b.Bytes(), nil
}

func (s *SalesReportResponse) ToEncoding(e encoding.Encoding) ([]byte, error) {
	switch e {
	case encoding.Json:
		return s.ToJson()
	case encoding.Csv:
		return s.ToCsv()
	}
	return nil, errors.New("I dont know how to encode that: " + e.String())
}
