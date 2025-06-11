package utils

import (
    "fmt"
    "log"
    "os"

    twilio "github.com/twilio/twilio-go"
    verify "github.com/twilio/twilio-go/rest/verify/v2"
)

func SendVerificationCode(phoneNumber string) error {
    accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
    authToken := os.Getenv("TWILIO_AUTH_TOKEN")
    serviceSid := os.Getenv("TWILIO_SERVICE_SID")

    if accountSid == "" || authToken == "" || serviceSid == "" {
        return fmt.Errorf("no twilio credentials")
    }

    client := twilio.NewRestClientWithParams(twilio.ClientParams{
        Username: accountSid,
        Password: authToken,
    })

    params := &verify.CreateVerificationParams{}
    params.SetTo(phoneNumber)
    params.SetChannel("sms")

    _, err := client.VerifyV2.CreateVerification(serviceSid, params)
    if err != nil {
        log.Printf("Twilio API Error: %+v\n", err)
        return fmt.Errorf("no se pudo enviar el código de verificación: %v", err)
    }

    return nil
}

func VerifyNumber(phoneNumber string, smsCode string) (*string, error) {
    accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
    authToken := os.Getenv("TWILIO_AUTH_TOKEN")
    serviceSid := os.Getenv("TWILIO_SERVICE_SID")

    if accountSid == "" || authToken == "" || serviceSid == "" {
        return nil, fmt.Errorf("no twilio credentials")
    }

    client := twilio.NewRestClientWithParams(twilio.ClientParams{
        Username: accountSid,
        Password: authToken,
    })

    params := &verify.CreateVerificationCheckParams{}
    params.SetTo(phoneNumber)
    params.SetCode(smsCode)

    verification, err := client.VerifyV2.CreateVerificationCheck(serviceSid, params)
    if err != nil {
        log.Printf("Twilio API Error Details: %+v\n", err)
        return nil, fmt.Errorf("failed to verify code: %v", err)
    }

    if verification.Status == nil {
        return nil, fmt.Errorf("no verification status received")
    }

    return verification.Status, nil
}
