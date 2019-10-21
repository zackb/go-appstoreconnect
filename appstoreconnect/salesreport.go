package appstoreconnect

func NewSalesReport(frequency string) *service {
	return &service{
		Path: "salesReports",
		Params: map[string]string{
			"filter[frequency]":     frequency,
			"filter[reportDate]":    "2019-09",
			"filter[reportType]":    "SALES",
			"filter[reportSubType]": "SUMMARY",
		},
	}
}
