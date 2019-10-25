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
	Reports []*SalesReportItem
}
type SalesReportItem struct {
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

type SalesReport service

const (
	Path = "salesReports"
)

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

func (f *Frequency) String() string {
	return string(*f)
}

func (r *ReportType) String() string {
	return string(*r)
}

func (r *ReportSubType) String() string {
	return string(*r)
}

// Clone deep copy
func (s *SalesReportItem) Clone() *SalesReportItem {
	c := SalesReportItem{}
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

// GetRange TODO
func (c *SalesReport) GetRange(timeRange *TimeRange, reportType ReportType, reportSubType ReportSubType) (*SalesReportResponse, error) {
	sr := SalesReportResponse{}
	for timeRange.Next() {
		s1, err := c.Get(timeRange.Current(), timeRange.Frequency, reportType, reportSubType)
		if err != nil {
			if err == ErrNoData {
				continue
			}
			return &sr, err
		}
		sr.Reports = append(sr.Reports, s1.Reports...)
	}
	return &sr, nil
}

func (c *SalesReport) Get(date time.Time, frequency Frequency, reportType ReportType, reportSubType ReportSubType) (*SalesReportResponse, error) {
	path := "salesReports"
	params := map[string]string{
		"filter[frequency]":     frequency.String(),
		"filter[reportDate]":    timeToReportDate(date, frequency),
		"filter[reportType]":    reportType.String(),
		"filter[reportSubType]": reportSubType.String(),
	}
	b, err := c.client.get(path, params)
	if err != nil {
		return nil, err
	}

	data := SalesReportItem{}
	p, err := encoding.NewTsvParser(bytes.NewReader(b), &data)
	if err != nil {
		return nil, err
	}

	ret := SalesReportResponse{}
	for {
		eof, err := p.Next()
		if eof {
			break
		}
		if err != nil {
			return nil, err
		}
		ret.Reports = append(ret.Reports, data.Clone())
	}
	return &ret, nil
}

func (s *SalesReportItem) GetHeader() []string {
	t := reflect.TypeOf(SalesReportItem{})
	tags := []string{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("tsv")
		tags = append(tags, tag)
	}

	return tags
}

func (s *SalesReportItem) Values() []string {
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
	d := SalesReportItem{}
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
