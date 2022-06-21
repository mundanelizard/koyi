package models

// PhoneNumber is based on the spec
// Reference https://en.wikipedia.org/wiki/E.164
type PhoneNumber struct {
	CountryCode      string `json:"countryCode" bson:"countryCode"`           // min 1 max 12
	SubscriberNumber string `json:"subscriberNumber" bson:"subscriberNumber"` // max of 12 digits
}

func (ph *PhoneNumber) IsValid() bool {
	if !isValidCountryCode(ph.CountryCode) {
		return false
	}

	if !isValidSubscriberNumber(ph.SubscriberNumber) {
		return false
	}

	return true
}

// isValidSubscriberNumber validates subscriber ar https://en.wikipedia.org/wiki/E.164
func isValidSubscriberNumber(sn string) bool {
	length := len(sn)
	return length > 3 && length <= 13
}

func isValidCountryCode(cc string) bool {
	length := len(cc)
	return length > 1 && length <= 3
}
