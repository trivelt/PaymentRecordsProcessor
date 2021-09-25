package paymentrecords

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
