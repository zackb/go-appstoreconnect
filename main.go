package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

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
	timeRange       *appstoreconnect.TimeRange
}

func main() {
	c, err := parseCmd()
	if err != nil {
		panic(err)
		flag.Usage()
		return
	}
	client, err := appstoreconnect.NewClientFromCredentialsFile(c.credentialsFile)
	checkError(err)
	e, err := c.execute(client)
	checkError(err)
	b, err := e.ToEncoding(c.outputFormat)
	checkError(err)
	fmt.Println(string(b))
}

func (c *cmd) execute(client *appstoreconnect.Client) (encoding.Encodable, error) {
	switch c.service {
	case CmdSalesReport:
		return client.SalesReport.GetRange(
			c.timeRange,
			appstoreconnect.ReportSales,
			appstoreconnect.SubReportSummary)
	case CmdFinanceReport:
		return client.FinanceReport.Get(c.timeRange.Start, "US")
	default:
		flag.Usage()
	}
	return nil, errors.New("No such command")
}

func parseCmd() (*cmd, error) {
	c := cmd{}
	var d string

	if len(os.Args) < 2 {
		return nil, errors.New("Must specify a command")
	}

	c.service = os.Args[1]

	fs := flag.NewFlagSet(c.service, flag.ExitOnError)
	fs.StringVar(&c.credentialsFile, "c", "credentials.yml", "path to credentials yaml file")
	fs.StringVar(&d, "d", "", "date string")
	fs.Var(&c.outputFormat, "o", "output format")
	fs.Parse(os.Args[2:])

	// default to json
	if c.outputFormat == encoding.None {
		c.outputFormat = encoding.Json
	}

	var err error
	if d == "" {
		c.timeRange = appstoreconnect.Yesterday()
	} else {
		c.timeRange, err = appstoreconnect.ParseTimeRange(d)

	}
	return &c, err
}

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}
