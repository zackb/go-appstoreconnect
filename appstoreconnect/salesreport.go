package appstoreconnect

import (
	"reflect"
	"time"
)

// Sales and Trends reports
// https://developer.apple.com/documentation/appstoreconnectapi/download_sales_and_trends_reports

type Frequency string

type ReportType string

type ReportSubType string

type SalesReportResponse struct {
	Provider              string `tsv:Provider`
	ProviderCountry       string `tsv:Provider Country`
	SKU                   string `tsv:SKU`
	Developer             string `tsv:Developer`
	Title                 string `tsv:Title`
	Version               string `tsv:Version`
	ProductTypeIdentifier string `tsv Product Type Identifier`
	Units                 int    `tsv:Units`
	DeveloperProceeds     string `tsv:Developer Proceeds`
	BeginDate             string `tsv:Begin Date`
	EndDate               string `tsv:End Date`
	CustomerCurrency      string `tsv:Customer Currency`
	CountryCode           string `tsv:Country Code`
	CurrencyOfProceeds    string `tsv:Currency of Proceeds`
	AppleIdentifier       string `tsv:Apple Identifier`
	CustomerPrice         string `tsv:Customer Price`
	PromoCode             string `tsv:Promo Code`
	ParentIdentifier      string `tsv:Parent Identifier`
	Subscription          string `tsv:Subscription`
	Period                string `tsv:Period`
	Category              string `tsv:Category`
	CMB                   string `tsv:CMB`
	Device                string `tsv:Device`
	SupportedPlatforms    string `tsv:Supported Platforms`
	ProceedsReason        string `tsv:Proceeds Reason`
	PreservedPricing      string `tsv:Preserved Pricing`
	Client                string `tsv:Client`
	OrderType             string `tsv:Order Type`
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
func (s *SalesReportResponse) Clone() *SalesReportResponse {
	c := SalesReportResponse{}
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

func (s *SalesReportResponse) GetHeader() []string {
	t := reflect.TypeOf(SalesReportResponse{})
	tags := []string{}
	for i := 0; i < t.NumField(); i++ {
		// Get the field, returns https://golang.org/pkg/reflect/#StructField
		field := t.Field(i)
		tag := field.Tag.Get("tsv")
		tags = append(tags, tag)

	}

	return tags
}
