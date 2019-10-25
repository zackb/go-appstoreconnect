package main

import (
	"errors"
	"flag"
	"fmt"
	"time"

	"github.com/zackb/go-appstoreconnect/appstoreconnect"
	"github.com/zackb/go-appstoreconnect/encoding"
)

const (
	CmdSalesReport   string = "SalesReport"
	CmdFinanceReport string = "FinanceReport"
)

type cmd struct {
	service         string
	credentialsFile string
	outputFormat    encoding.Encoding
	dateRange       string
}

func main() {
	c, err := parseCmd()
	if err != nil {
		flag.Usage()
		return
	}
	client, err := appstoreconnect.NewClientFromCredentialsFile(c.credentialsFile)
	checkError(err)
	e, err := c.execute(client)
	checkError(err)
	fmt.Println(e.ToEncoding(c.outputFormat))
}

func (c *cmd) execute(client *appstoreconnect.Client) (encoding.Encodable, error) {
	switch c.service {
	case CmdSalesReport:
		return client.SalesReport.GetRange(
			appstoreconnect.NewTimeRange(
				appstoreconnect.NewTime("2019-07-27"),
				appstoreconnect.NewTime("2019-10-21"),
				appstoreconnect.Monthly,
			),
			appstoreconnect.ReportSales,
			appstoreconnect.SubReportSummary)
	case CmdFinanceReport:
		return client.FinanceReport.Get(time.Now(), "US")
	default:
		flag.Usage()
	}
	return nil, errors.New("No such command")
}

func parseCmd() (*cmd, error) {
	c := cmd{}
	flag.StringVar(&c.credentialsFile, "c", "credentials.yml", "path to credentials yaml file")
	flag.StringVar(&c.dateRange, "d", "", "date string")
	flag.Var(&c.outputFormat, "o", "output format")
	c.service = flag.Arg(0)
	if c.service == "" {
		return nil, errors.New("Must specify a command")
	}

	if c.dateRange == "" {
		c.dateRange = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	}
	return &c, nil
}

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}
