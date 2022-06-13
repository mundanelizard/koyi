package helpers

type Sms struct {
	Text *string
}

func (s *Sms) Send() error {
	return nil
}

type Sendable interface {
	Send() error
}
