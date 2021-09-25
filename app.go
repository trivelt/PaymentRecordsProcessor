package main

import (
	"github.com/trivelt/payment-records-processor/paymentrecords"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("app")

func setupLogging(verbose bool) {
	var format = logging.MustStringFormatter(
		`%{color} %{level} %{color:reset} %{message}`,
	)
	if verbose {
		logging.SetLevel(logging.DEBUG, "")
	} else {
		logging.SetLevel(logging.INFO, "")
	}
	logging.SetFormatter(format)
}

func main() {
	config := paymentrecords.ParseArgs()
	setupLogging(config.Verbose)
	data := paymentrecords.FetchData()
	records := paymentrecords.Transform(data)
	paymentrecords.Upload(records, config.ApiUrl, config.ApiKey, config.SingleRequestUpload)
}
