# eCatholic Coding Practical

Application for simple processing of payment records. It fetches CSV data from the network, performs some 
transformations and then uploads data into the specified URL. The program is written in Go. 

## Table of Content
- [eCatholic Coding Practical](#ecatholic-coding-practical)
  * [Author](#author)
  * [Introduction](#introduction)
  * [Usage](#usage)
  * [Upload requests](#upload-requests)
  * [Tests](#tests)
  * [Implementation](#implementation)
  * [Possible improvements](#possible-improvements)
  * [Acknowledgements](#acknowledgements)

## Author
**Maciej Michalec** - polyglot developer (creator of [PolyDev.PL](https://polydev.pl)), with knowledge of Golang but without prior commercial experience in this language. 

## Introduction
The application follows the provided specification and implements the following actions:
1. Downloads the structured data file from https://s3.amazonaws.com/ecatholic-hiring/data.csv.
2. Removes the field named 'Memo' from all records.
3. Adds a field named "import_date" and populate it appropriately.
4. For any record that has an empty value, sets the value of the field to the value "missing".
5. Removes any duplicate records.
6. Submits the records as JSON objects named `PaymentRecord` to a REST API url provided by the user. 
`POST` method is used.

## Usage

In order to run the application, please download the repository and execute:
```
go run app.go -url "http://UPLOAD_DATA_ENDPOINT.com"
``` 

All required Go packages should be automatically fetched.

The application accepts the following command-line arguments:
 * `-key API_KEY` - API Key used to authorize the upload request (used in the `X-API-KEY` header).
 * `-single-upload` - flag which indicate that the app should upload all payment records in a single request instead of the separate request for each record.
 * `-url "UPLOAD_DATA_URL"` - URL of the API endpoint used to upload the data (it should accept POST method). [**Required**]
 * `-verbose` - flag which enables verbose mode (more information printed on the output).

If you don't have your own server accepting POST requests sent by the application, you can use some 
online requests bin, like https://hookbin.com, to verify the data.  

## Upload requests
When using a `-single-upload` flag, the upload request contains a list of the payment records:

```json
{
  "PaymentRecord": [
    {
      "Date": "01/04/2016",
      "Name": "Emil Holder",
      "Address": "4040 Melville Street",
      "Address2": "missing",
      "City": "Arlington",
      "State": "TN",
      "Zipcode": "38002",
      "Telephone": "missing",
      "Mobile": "731-513-2214",
      "Amount": "$50",
      "Processor": "PayPal",
      "import_date": "2021-09-26"
    },
    {
      "Date": "06/15/2019",
      "Name": "Mike Smith",
      "Address": "2483 Farland Avenue",
      "Address2": "missing",
      "City": "Warrensburg",
      "State": "MO",
      "Zipcode": "64093",
      "Telephone": "443-323-6215",
      "Mobile": "410-726-6477",
      "Amount": "$40",
      "Processor": "Stripe",
      "import_date": "2021-09-26"
    }
  ]
}
```

Without this flag, the application sends separate request for each record:

```json
{
  "PaymentRecord": {
    "Date": "05/23/2019",
    "Name": "Shirley Smith",
    "Address": "2451 Fairfax Drive",
    "Address2": "missing",
    "City": "Newark",
    "State": "NJ",
    "Zipcode": "7102",
    "Telephone": "320-763-7283",
    "Mobile": "missing",
    "Amount": "$40",
    "Processor": "PayPal",
    "import_date": "2021-09-26"
  }
}
```

## Tests
File `paymentrecords/paymentrecords_test.go` contains a small set of the automated tests
You can run these tests by executing the following command in the root directory of the application:

```
go test ./... -v 
``` 

## Implementation
When creating software, I primarily try to make the code as readable and maintainable as possible. 
Thus is divided it into lots of small functions, s whose workings should be easy to understand. 

I also modularized the application, moving the whole business logic into the separate package (`github.com/trivelt/payment-records-processor/paymentrecords`) and dividing
it into multiple short files.

Reading CSV files and sending JSON files was implemented using 
an approach popular in the Go language that uses data structures with tags, used automatically during (un)marshaling.    

## Possible improvements
 * I am aware that the code coverage is not the best and having more time for completing this exercise I would certainly add more tests
 * Dockerization of the app would be a nice thing, especially if the program were to be extended
 * More complex configuration, allowing e.g. for specifying the fetch URL and using authorization to get the data.

## Acknowledgements
Thank you for reading this README file and reviewing my solution. 
