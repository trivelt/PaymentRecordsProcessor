package main

import (
	"flag"
	"fmt"
	"os"
	"bytes"
	"time"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/gocarina/gocsv"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("app")


type Set struct {
	list map[PaymentRecord]struct{}
}

func (set *Set) Add(record PaymentRecord) {
	set.list[record] = struct{}{}
}

func (set *Set) List() []PaymentRecord {
    records := make([]PaymentRecord, 0, len(set.list))
    for key := range set.list {
        records = append(records, key)
    }
	return records
}

func NewSet() *Set {
	set := &Set{}
	set.list = make(map[PaymentRecord]struct{})
	return set
}


type PaymentRecord struct {
	Date		string `csv:"default=missing"`
	Name		string `csv:"default=missing"`
	Address		string `csv:"default=missing"`
	Address2	string `csv:"default=missing"`
	City		string `csv:"default=missing"`
	State		string `csv:"default=missing"`
	Zipcode		string `csv:"default=missing"`
	Telephone	string `csv:"default=missing"`
	Mobile		string `csv:"default=missing"`
	Amount		string `csv:"default=missing"`
	Processor	string `csv:"default=missing"`
	ImportDate	string `json:"import_date"`
}

type MultiplePaymentRecordsMessage struct {
	Data []PaymentRecord `json:"PaymentRecord"`
}

type SinglePaymentRecordMessage struct {
	Data PaymentRecord `json:"PaymentRecord"`
}


type AppConfig struct {
	ApiUrl				string
	ApiKey				string
	Verbose				bool
	SingleRequestUpload	bool
}


func addImportDate(records []*PaymentRecord) {
	currentTime := time.Now().Local()
	importDate := currentTime.Format("2006-01-02")

	for _, record := range records {
		log.Debugf("Fetched record %s", *record)
		record.ImportDate = importDate
	}
	log.Debugf("Added import date %s", importDate)
}

func fetchData() string {
	resp, err := http.Get("https://s3.amazonaws.com/ecatholic-hiring/data.csv")
	if err != nil {
		log.Errorf("Cannot fetch CSV data: %s", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

    data, err := ioutil.ReadAll(resp.Body)
    if err != nil {
		panic(err)
	}

	return string(data)
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

func isErrorResponse(resp *http.Response) bool {
	return resp.StatusCode >= 300
}

func uploadData(requestJson []byte, url string, apiKey string) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestJson))
	if err != nil {
		panic(err)
	}
	req.Header.Add("X-API-KEY", apiKey)
	req.Header.Add("Content-Type", "application/json")

	client := http.Client{Timeout: time.Duration(1) * time.Second}
	resp, err := client.Do(req)

	 if err != nil {
		log.Errorf("Cannot upload Payment Records: %s", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if isErrorResponse(resp) {
		log.Errorf("Received error response when uploading Payment Records: %s", resp.Status)
		os.Exit(1)
	}
}

func sentRequestWithRecordsList(records Set, url string, apiKey string) { 
	log.Debugf("All Payment Records to upload: %s", records.list)
	message := MultiplePaymentRecordsMessage{Data: records.List()}
	requestJson, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}

	uploadData(requestJson, url, apiKey)
	log.Info("All Payment Records submitted successfully")
}

func sentRequestForEachRecord(records Set, url string, apiKey string) {
	for record, _ := range records.list {
		log.Debugf("Uploading Payment Record %s", record)

		message := SinglePaymentRecordMessage{Data: record}
		requestJson, err := json.Marshal(message)
		if err != nil {
			panic(err)
		}
	
		uploadData(requestJson, url, apiKey)
	}
	log.Info("All Payment Records submitted successfully")
}


func upload(records Set, url string, apiKey string, singleRequestUpload bool) {
	if singleRequestUpload {
		sentRequestWithRecordsList(records, url, apiKey)
	} else {
		sentRequestForEachRecord(records, url, apiKey)
	}
}


func parseArgs() AppConfig {
	apiUrl := flag.String("url", "", "URL of the API endpoint used to upload the data (it should accept POST method). [Required]")
	apiKey := flag.String("key", "", "API Key used to authorize the upload request. [Required]")
	verbose := flag.Bool("verbose", false, "Enable verbose mode.")
	singleRequestUpload := flag.Bool("single-upload", false, "Upload all payment records in a single request instead of the separate request for each record.")

    flag.Parse()
    if *apiUrl == "" || *apiKey == ""{
		fmt.Println("Please provide all required parameters:")
        flag.PrintDefaults()
        os.Exit(1)
    }

	return AppConfig{ApiUrl: *apiUrl, ApiKey: *apiKey, Verbose: *verbose, SingleRequestUpload: *singleRequestUpload}
}

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
	config := parseArgs()
	setupLogging(config.Verbose)
	body := fetchData()
	records := unpackToCsv(body)
	addImportDate(records)
	uniqueRecords := makeUnique(records)
	upload(uniqueRecords, config.ApiUrl, config.ApiKey, config.SingleRequestUpload)
}
