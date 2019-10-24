package appstoreconnect

import "time"

// Finance Report
// https://developer.apple.com/documentation/appstoreconnectapi/download_finance_reports

const (
	PathFinanceReport = "financeReports"
)

type FinanceReport struct {
	service
}

func (f *FinanceReport) Get(date time.Time, regionCode string) ([]byte, error) {

	return f.client.get(PathFinanceReport,
		map[string]string{
			"filter[regionCode]": regionCode,
			"filter[reportDate]": timeToReportDate(date, Monthly),
			"filter[reportType]": "FINANCIAL",
		},
	)
}
