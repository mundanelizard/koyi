package services

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/mundanelizard/koyi/server/config"
)

const (
	CharSet = "UTF-8"
)

type email struct {
	To      string
	Subject string
	Text    *string
	HTML    *string
}

var sess, err = session.NewSession(&aws.Config{
	Region: aws.String("us-west-2")}, // update to the right one.
)

// Create an SES session.
var svc = ses.New(sess)

func NewMail(to, subject string, text, html *string) Sendable {
	return &email{
		to,
		subject,
		text,
		html,
	}
}

func (email *email) Send() error {
	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(email.To),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(*email.HTML),
				},
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(*email.Text),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(email.Subject),
			},
		},
		Source: aws.String(config.EmailAddress),
		// Uncomment to use a configuration set
		//ConfigurationSetName: aws.String(ConfigurationSet),
	}

	// Attempt to send the email.
	_, err = svc.SendEmail(input)

	if err != nil {
		return err
	}

	return nil
}
