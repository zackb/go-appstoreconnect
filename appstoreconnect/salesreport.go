package appstoreconnect

// Sales and Trends reports
// https://developer.apple.com/documentation/appstoreconnectapi/download_sales_and_trends_reports

type Frequency string

type ReportType string

type ReportSubType string

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

func NewSalesReport(frequency Frequency, reportType ReportType, reportSubType ReportSubType) *service {
	return &service{
		Path: "salesReports",
		Params: map[string]string{
			"filter[frequency]":     frequency.String(),
			"filter[reportDate]":    "2019-09",
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
