package appstoreconnect

import (
	"time"

	"github.com/zackb/go-appstoreconnect/encoding"
)

// Finance Report
// https://developer.apple.com/documentation/appstoreconnectapi/download_finance_reports

const (
	pathFinanceReports = "financeReports"
)

type FinanceReport service

type FinanceReportResponse struct {
	raw []byte
}

func (f *FinanceReport) Get(date time.Time, regionCode string) (*FinanceReportResponse, error) {

	b, err := f.client.get(pathFinanceReports,
		map[string]string{
			"filter[regionCode]": regionCode,
			"filter[reportDate]": timeToReportDate(date, Monthly),
			"filter[reportType]": "FINANCIAL",
		},
	)

	return &FinanceReportResponse{raw: b}, err
}

// I only make free apps so I dont know what this returns!
func (f *FinanceReportResponse) ToEncoding(e encoding.Encoding) ([]byte, error) {
	return f.raw, nil
}
