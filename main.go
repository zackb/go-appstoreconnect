package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/zackb/go-appstoreconnect/appstoreconnect"
)

type cmd struct {
	credentialsFile string
	outputFormat    Encoding
}

func main() {
	c, _ := parseCmd()
	fmt.Println(c)
	client, err := appstoreconnect.NewClientFromCredentialsFile(c.credentialsFile)
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

func parseCmd() (*cmd, error) {
	c := cmd{}
	flag.StringVar(&c.credentialsFile, "c", "credentials.yml", "path to credentials yaml file")
	flag.Var(&c.outputFormat, "o", "output format")

	flag.Parse()
	return &c, nil
}
