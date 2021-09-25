package paymentrecords

import (
	"os"
	"bytes"
	"time"
	"net/http"
	"encoding/json"
)

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
	log.Debugf("All Payment Records to upload: %s", records.data)
	message := MultiplePaymentRecordsMessage{Data: records.List()}
	requestJson, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}

	uploadData(requestJson, url, apiKey)
	log.Info("All Payment Records submitted successfully")
}

func sentRequestForEachRecord(records Set, url string, apiKey string) {
	for record, _ := range records.data {
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


func Upload(records Set, url string, apiKey string, singleRequestUpload bool) {
	if singleRequestUpload {
		sentRequestWithRecordsList(records, url, apiKey)
	} else {
		sentRequestForEachRecord(records, url, apiKey)
	}
}
