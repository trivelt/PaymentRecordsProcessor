package paymentrecords

import (
	"os"
	"time"
	"github.com/gocarina/gocsv"
	"github.com/op/go-logging"
)


var log = logging.MustGetLogger("app")


func addImportDate(records []*PaymentRecord) {
	currentTime := time.Now().Local()
	importDate := currentTime.Format("2006-01-02")

	for _, record := range records {
		log.Debugf("Fetched record %s", *record)
		record.ImportDate = importDate
	}
	log.Debugf("Added import date %s", importDate)
}


func unpackToCsv(body string) (records []*PaymentRecord) {
	if err := gocsv.UnmarshalString(body, &records); err != nil {
		log.Errorf("Cannot unpack received CSV data: %s", err)
		os.Exit(1)
	}
	return records
}

func makeUnique(records []*PaymentRecord) Set {
	uniqueRecords := NewSet()
	for _, record := range records {
		uniqueRecords.Add(*record)
	}
	return *uniqueRecords
}

func Transform(body string) Set {
	records := unpackToCsv(body)
	addImportDate(records)
	return makeUnique(records)
}
