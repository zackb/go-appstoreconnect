package main

import (
	"fmt"
	"time"

	"github.com/zackb/go-appstoreconnect/appstoreconnect"
)

func main() {
	client, err := appstoreconnect.NewClientFromCredentialsFile("credentials.yml")
	if err != nil {
		panic(err)
	}

	b, err := client.GetSalesReport(
		time.Now(),
		appstoreconnect.Weekly,
		appstoreconnect.ReportSales,
		appstoreconnect.SubReportSummary)

	if err != nil {
		panic(err)
	}
	for _, r := range b {
		fmt.Println(r.Units)
	}
}
