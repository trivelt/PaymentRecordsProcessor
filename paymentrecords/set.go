package paymentrecords

type Set struct {
	data map[PaymentRecord]struct{}
}

func (set *Set) Add(record PaymentRecord) {
	set.data[record] = struct{}{}
}

func (set *Set) List() []PaymentRecord {
    records := make([]PaymentRecord, 0, len(set.data))
    for key := range set.data {
        records = append(records, key)
    }
	return records
}

func NewSet() *Set {
	set := &Set{}
	set.data = make(map[PaymentRecord]struct{})
	return set
}
