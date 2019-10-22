package appstoreconnect

import "time"

// Finance Report
// https://developer.apple.com/documentation/appstoreconnectapi/download_finance_reports

func NewFinanceReport(date time.Time, regionCode string) *service {
	return &service{
		Path: "financeReports",
		Params: map[string]string{
			"filter[regionCode]": regionCode,
			"filter[reportDate]": timeToReportDate(date, Monthly),
			"filter[reportType]": "FINANCIAL",
		},
	}
}
