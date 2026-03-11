package appstoreconnect

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/zackb/go-appstoreconnect/encoding"
)

// Finance Report
// https://developer.apple.com/documentation/appstoreconnectapi/download_finance_reports

const (
	pathFinanceReports = "financeReports"
)

// FinanceReport service is responsible for communicating with the "financeReports" endpoint
type FinanceReport service

// FinanceReportResponse represents the response from a finance report request
type FinanceReportResponse struct {
	Reports []*FinanceReportItem
}

// FinanceReportItem is information returned about FINANCIAL reports
type FinanceReportItem struct {
	StartDate                        string  `tsv:"Start Date"`
	EndDate                          string  `tsv:"End Date"`
	UPC                              string  `tsv:"UPC"`
	ISRCISBN                         string  `tsv:"ISRC/ISBN"`
	VendorIdentifier                 string  `tsv:"Vendor Identifier"`
	Quantity                         string  `tsv:"Quantity"`
	PartnerShare                     string  `tsv:"Partner Share"`
	ExtendedPartnerShare             string  `tsv:"Extended Partner Share"`
	PartnerShareCurrency             string  `tsv:"Partner Share Currency"`
	SalesOrReturn                    string  `tsv:"Sales or Return"`
	AppleIdentifier                  string  `tsv:"Apple Identifier"`
	ArtistShowDeveloperAuthor        string  `tsv:"Artist/Show/Developer/Author"`
	Title                            string  `tsv:"Title"`
	LabelStudioNetworkDeveloperPublisher string `tsv:"Label/Studio/Network/Developer/Publisher"`
	Grid                             string  `tsv:"Grid"`
	ProductTypeIdentifier            string  `tsv:"Product Type Identifier"`
	ISANOtherIdentifier              string  `tsv:"ISAN/Other Identifier"`
	CountryOfSale                    string  `tsv:"Country Of Sale"`
	PreOrderFlag                     string  `tsv:"Pre-order Flag"`
	PromoCode                        string  `tsv:"Promo Code"`
	CustomerPrice                    string  `tsv:"Customer Price"`
	CustomerCurrency                 string  `tsv:"Customer Currency"`
}

// Clone deep copy
func (f *FinanceReportItem) Clone() *FinanceReportItem {
	c := *f
	return &c
}

func (f *FinanceReportItem) GetHeader() []string {
	t := reflect.TypeFor[FinanceReportItem]()
	tags := []string{}
	for field := range t.Fields() {
		tag := field.Tag.Get("tsv")
		tags = append(tags, tag)
	}
	return tags
}

func (f *FinanceReportItem) Values() []string {
	vals := []string{}
	t := reflect.ValueOf(f).Elem()
	for i := range t.NumField() {
		vals = append(vals, t.Field(i).String())
	}
	return vals
}

// Get fetches a finance report for a given month and region code.
// Finance reports are always MONTHLY, so only the year and month of date are used.
// regionCode is the two-letter region code (e.g. "US", "ZZ" for worldwide).
func (f *FinanceReport) Get(date time.Time, regionCode string) (*FinanceReportResponse, error) {
	b, err := f.client.get(pathFinanceReports,
		map[string]string{
			"filter[regionCode]": regionCode,
			"filter[reportDate]": timeToReportDate(date, Monthly),
			"filter[reportType]": "FINANCIAL",
		},
	)
	if err != nil {
		return nil, err
	}

	item := FinanceReportItem{}
	p, err := encoding.NewTsvParser(bytes.NewReader(stripFinanceFooter(b)), &item)
	if err != nil {
		return nil, err
	}

	ret := FinanceReportResponse{}
	for {
		eof, err := p.Next()
		if eof {
			break
		}
		if err != nil {
			return nil, err
		}
		ret.Reports = append(ret.Reports, item.Clone())
	}
	return &ret, nil
}

// GetRange fetches finance reports for every month in the given TimeRange.
// ErrNoData months are silently skipped.
func (f *FinanceReport) GetRange(tr *TimeRange, regionCode string) (*FinanceReportResponse, error) {
	ret := FinanceReportResponse{}
	for tr.Next() {
		r, err := f.Get(tr.Current(), regionCode)
		if err != nil {
			if err == ErrNoData {
				continue
			}
			return &ret, err
		}
		ret.Reports = append(ret.Reports, r.Reports...)
	}
	return &ret, nil
}

// stripFinanceFooter removes Apple's trailing summary rows (Total_Rows, Total_Amount,
// Total_Units) which have fewer fields than the header and break the csv.Reader.
func stripFinanceFooter(b []byte) []byte {
	var out bytes.Buffer
	sc := bufio.NewScanner(bytes.NewReader(b))
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "Total_") {
			continue
		}
		out.WriteString(line)
		out.WriteByte('\n')
	}
	return out.Bytes()
}

func (f *FinanceReportResponse) ToJson() ([]byte, error) {
	return json.Marshal(f)
}

func (f *FinanceReportResponse) ToCsv() ([]byte, error) {
	b := bytes.Buffer{}
	buf := bufio.NewWriter(&b)

	w := csv.NewWriter(buf)
	d := FinanceReportItem{}
	headers := d.GetHeader()
	if err := w.Write(headers); err != nil {
		return nil, err
	}

	for _, r := range f.Reports {
		values := r.Values()
		if err := w.Write(values); err != nil {
			return nil, err
		}
	}
	buf.Flush()
	return b.Bytes(), nil
}

func (f *FinanceReportResponse) ToTsv() ([]byte, error) {
	b := bytes.Buffer{}
	buf := bufio.NewWriter(&b)

	w := csv.NewWriter(buf)
	w.Comma = '\t'
	d := FinanceReportItem{}
	headers := d.GetHeader()
	if err := w.Write(headers); err != nil {
		return nil, err
	}

	for _, r := range f.Reports {
		values := r.Values()
		if err := w.Write(values); err != nil {
			return nil, err
		}
	}
	buf.Flush()
	return b.Bytes(), nil
}

func (f *FinanceReportResponse) ToEncoding(e encoding.Encoding) ([]byte, error) {
	switch e {
	case encoding.Json:
		return f.ToJson()
	case encoding.Csv:
		return f.ToCsv()
	case encoding.Tsv:
		return f.ToTsv()
	}
	return nil, errors.New("I dont know how to encode that: " + e.String())
}
