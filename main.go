package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/zackb/go-appstoreconnect/appstoreconnect"
	"github.com/zackb/go-appstoreconnect/encoding"
)

type cmd struct {
	credentialsFile string
	outputFormat    encoding.Encoding
}

func main() {
	c, _ := parseCmd()
	fmt.Println(c)
	client, err := appstoreconnect.NewClientFromCredentialsFile(c.credentialsFile)
	if err != nil {
		panic(err)
	}

	d, err := client.GetSalesReport(
		time.Now(),
		appstoreconnect.Weekly,
		appstoreconnect.ReportSales,
		appstoreconnect.SubReportSummary)

	if err != nil {
		panic(err)
	}

	b, err := encoding.NewJsonEncoder().Encode(d)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}

func parseCmd() (*cmd, error) {
	c := cmd{}
	flag.StringVar(&c.credentialsFile, "c", "credentials.yml", "path to credentials yaml file")
	flag.Var(&c.outputFormat, "o", "output format")

	flag.Parse()
	return &c, nil
}
