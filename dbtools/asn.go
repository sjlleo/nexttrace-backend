package dbtools

type Asn struct {
	Asn     int
	Alias   string
	Isp     string
	Website string
}

func SearchASNData(asnNumber string) (*Asn, error) {
	// ASN 数据
	db := GetDB()
	asn := Asn{}
	res := db.Where("asn = ?", asnNumber).Find(&asn)
	return &asn, res.Error
}
