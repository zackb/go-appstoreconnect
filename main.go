package main

import (
	"fmt"
	"github.com/zackb/go-appstoreconnect/appstoreconnect"
)


func main() {
	client, err := appstoreconnect.NewClientFromCredentialsFile("credentials.yml")
	if err != nil {
		panic(err)
	}

	b, err := client.GetSalesReport(
		appstoreconnect.Daily,
		appstoreconnect.ReportSales,
		appstoreconnect.SubReportSummary)

	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
}
