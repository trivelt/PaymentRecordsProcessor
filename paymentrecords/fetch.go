package paymentrecords

import (
	"os"
	"net/http"
	"io/ioutil"
)

const CSV_FILE_URL string = "https://s3.amazonaws.com/ecatholic-hiring/data.csv"


func FetchData() string {
	resp, err := http.Get(CSV_FILE_URL)
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
