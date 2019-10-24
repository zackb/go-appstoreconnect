package main

import (
	"flag"
	"fmt"

	"github.com/zackb/go-appstoreconnect/appstoreconnect"
	"github.com/zackb/go-appstoreconnect/encoding"
)

type cmd struct {
	credentialsFile string
	outputFormat    encoding.Encoding
}

func main() {
	c, _ := parseCmd()
	client, err := appstoreconnect.NewClientFromCredentialsFile(c.credentialsFile)
	checkError(err)

	// financeReport(client, c)
	salesReport(client, c)
}

/*
func financeReport(client *appstoreconnect.Client, c *cmd) {
	d, err := client.GetFinanceReport(time.Now(), "US")
	checkError(err)
	fmt.Println(string(d))
}
*/

func salesReport(client *appstoreconnect.Client, c *cmd) {
	/*
		d, err := client.GetSalesReport(
			time.Now().AddDate(0, -2, 0),
			appstoreconnect.Weekly,
			appstoreconnect.ReportSales,
			appstoreconnect.SubReportSummary)
	*/
	/*
		appstoreconnect.NewTimeRange(
			time.Now().Add(-time.Hour*24*3),
			time.Now().Add(-time.Hour*24),
			appstoreconnect.Daily,
		),
	*/
	// appstoreconnect.NewSingleTimeRange("2019-09-01", appstoreconnect.Monthly),
	d, err := client.SalesReport.GetRange(
		appstoreconnect.NewTimeRange(
			appstoreconnect.NewTime("2019-07-27"),
			appstoreconnect.NewTime("2019-10-21"),
			appstoreconnect.Monthly,
		),
		appstoreconnect.ReportSales,
		appstoreconnect.SubReportSummary)

	checkError(err)
	b, err := d.ToEncoding(c.outputFormat)
	// b, err := encoding.NewJsonEncoder().Encode(d)
	checkError(err)
	fmt.Println(string(b))
}

func parseCmd() (*cmd, error) {
	c := cmd{}
	flag.StringVar(&c.credentialsFile, "c", "credentials.yml", "path to credentials yaml file")
	flag.Var(&c.outputFormat, "o", "output format")

	flag.Parse()
	return &c, nil
}

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}
