package paymentrecords

import (
	"flag"
	"os"
	"fmt"
)

type AppConfig struct {
	ApiUrl				string
	ApiKey				string
	Verbose				bool
	SingleRequestUpload	bool
}

func ParseArgs() AppConfig {
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
