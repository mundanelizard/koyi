package email

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

const (
	CharSet = "UTF-8"
)

type Email struct {
	To       string
	From     string
	Subject  string
	BodyText *string
	BodyHTML *string
}

var sess, err = session.NewSession(&aws.Config{
	Region: aws.String("us-west-2")}, // update to the right one.
)

// Create an SES session.
var svc = ses.New(sess)

func (email *Email) Send() (bool, error) {
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
					Data:    aws.String(*email.BodyHTML),
				},
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(*email.BodyText),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(email.Subject),
			},
		},
		Source: aws.String(email.From),
		// Uncomment to use a configuration set
		//ConfigurationSetName: aws.String(ConfigurationSet),
	}

	// Attempt to send the email.
	_, err = svc.SendEmail(input)

	if err != nil {
		return false, err
	}

	return true, nil
}
