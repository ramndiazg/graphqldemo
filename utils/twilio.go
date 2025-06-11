package utils

import (
	"os"

	twilio "github.com/twilio/twilio-go"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
)

func VerifyNumber(phoneNumber string, smsCode string) (*string, error) {

	TWILIO_SERVICE_SID := os.Getenv("TWILIO_SERVICE_SID")

	client := twilio.NewRestClient()
	params := &verify.CreateVerificationCheckParams{}
	params.SetTo(phoneNumber)
	params.SetCode(smsCode)

	verification, verificationErr := client.VerifyV2.CreateVerificationCheck(string(TWILIO_SERVICE_SID), params)

	if verificationErr != nil {
		return nil, verificationErr
	}

	return verification.Status, nil
}
