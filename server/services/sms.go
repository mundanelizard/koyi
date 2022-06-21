package services

type sms struct {
	Text *string
}

func (s *sms) Send() error {
	return nil
}

type Sendable interface {
	Send() error
}

func NewSms(text *string) Sendable {
	return &sms{Text: text}
}
