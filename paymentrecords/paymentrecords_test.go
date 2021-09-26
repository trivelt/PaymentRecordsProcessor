package paymentrecords

import (
	"time"
	"testing"
	"net/http"
	"net/http/httptest"
	"github.com/stretchr/testify/assert"
)


func TestTransformData_ShouldUnpackCsvFile(t* testing.T) {
	data := `
Name,City,Telephone,Address,Address2,State,Mobile,Amount,Processor,Zipcode,Date
John,NY,1234-5678,SomeAddress,AddressNo2,TestState,01234-0042,USD 10,TestProc,ZIP-CODE,2021-09-26
`
	records := Transform(data)
	paymentRecord := records.List()[0]

	assert.Equal(t, "John", paymentRecord.Name)
	assert.Equal(t, "NY", paymentRecord.City)
	assert.Equal(t, "1234-5678", paymentRecord.Telephone)
	assert.Equal(t, "SomeAddress", paymentRecord.Address)
	assert.Equal(t, "AddressNo2", paymentRecord.Address2)
	assert.Equal(t, "TestState", paymentRecord.State)
	assert.Equal(t, "01234-0042", paymentRecord.Mobile)
	assert.Equal(t, "USD 10", paymentRecord.Amount)
	assert.Equal(t, "TestProc", paymentRecord.Processor)
	assert.Equal(t, "ZIP-CODE", paymentRecord.Zipcode)
	assert.Equal(t, "2021-09-26", paymentRecord.Date)

}

func TestTransformData_ShouldAddImportDate(t* testing.T) {
	data := "Name\nJohn\nBob"
	records := Transform(data)
	firstRecord := records.List()[0]
	secondRecord := records.List()[1]

	expectedDate := time.Now().Local().Format("2006-01-02")
	assert.Equal(t, expectedDate, firstRecord.ImportDate)
	assert.Equal(t, expectedDate, secondRecord.ImportDate)
}

func TestTransformData_ShouldStoreOnlyUniqueRecords(t* testing.T) {
	data := `
Name,City,Telephone
John,NY,123
John,NY,1234
John,NY,123
`

	records := Transform(data)
	assert.Equal(t, 2, len(records.List()))
}

func TestUpload_ShouldSendMultipleRequests(t *testing.T) {
	requestsCounter := 0
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		requestsCounter++
		rw.Write([]byte(`OK`))
	}))
	defer server.Close()


	records := NewSet()
	records.Add(PaymentRecord{Name: "FirstRecord"})
	records.Add(PaymentRecord{Name: "SecondRecord"})
	singleRequestUpload := false
	Upload(*records, server.URL, "someKey", singleRequestUpload)

	assert.Equal(t, requestsCounter, 2)
}

func TestUpload_ShouldSendSingleRequest(t *testing.T) {
	requestsCounter := 0
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		requestsCounter++
		rw.Write([]byte(`OK`))
	}))
	defer server.Close()


	records := NewSet()
	records.Add(PaymentRecord{Name: "FirstRecord"})
	records.Add(PaymentRecord{Name: "SecondRecord"})
	singleRequestUpload := true
	Upload(*records, server.URL, "someKey", singleRequestUpload)

	assert.Equal(t, requestsCounter, 1)
}

func TestUpload_ShouldUseProvidedApiKey(t *testing.T) {
	apiKey := "MY_API_KEY"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		assert.Equal(t, apiKey, req.Header.Get("X-API-KEY"))
		rw.Write([]byte(`OK`))
	}))
	defer server.Close()


	records := NewSet()
	records.Add(PaymentRecord{Name: "FirstRecord"})
	Upload(*records, server.URL, apiKey, false)
}
