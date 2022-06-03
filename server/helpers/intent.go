package helpers

import (
	"github.com/mundanelizard/koyi/server/models"
)

var actionSubject = map[string]string{
	models.VerifyEmailIntent:       "https://gmail.com/google.com",
	models.VerifyPhoneNumberIntent: "https://email/intent/intent",
}

var actionHTML = map[string]string{
	models.VerifyEmailIntent:       "https://gmail.com/google.com",
	models.VerifyPhoneNumberIntent: "https://email/intent/intent",
}

var actionText = map[string]string{
	models.VerifyEmailIntent:       "https://gmail.com/google.com",
	models.VerifyPhoneNumberIntent: "https://email/intent/intent",
}

func GetEmailDetails(action string) (string, *string, *string) {
	// todo => use a template file
	text := actionText[action]
	html := actionHTML[action]
	subject := actionSubject[action]

	return subject, &text, &html
}
