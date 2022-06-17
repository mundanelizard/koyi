package models

// PhoneNumber is based on the spec
// Reference https://en.wikipedia.org/wiki/E.164
type PhoneNumber struct {
	CountryCode      string `json:"countryCode" bson:"countryCode"`           // min 1 max 12
	SubscriberNumber string `json:"subscriberNumber" bson:"subscriberNumber"` // max of 12 digits
}
