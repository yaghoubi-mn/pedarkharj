package sms

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

type SMSInput struct {
	Mobile     string         `json:"mobile"`
	TemplateId int            `json:"templateId"`
	Parameters []SMSParameter `json:"parameters"`
}

type SMSParameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type SMSOutput struct {
	Status  int             `json:"status"`
	Message string          `json:"message"`
	Data    []SMSOutputData `json:"data"`
}

type SMSOutputData struct {
	MessageID int     `json:"messageId"`
	Cost      float32 `json:"cost"`
}

func SendOTPSMS(mobile string, otp int) error {
	fmt.Println(mobile, "------------------")
	templateId, err := strconv.Atoi(os.Getenv("SMS_TEMPLATE_ID"))
	if err != nil {
		return err
	}
	smsInput := SMSInput{
		Mobile:     mobile,
		TemplateId: templateId,
		Parameters: []SMSParameter{
			{
				Name:  "Code",
				Value: strconv.Itoa(otp),
			},
		},
	}
	jsonSMS, err := json.Marshal(smsInput)
	if err != nil {
		return err
	}
	request, err := http.NewRequest(http.MethodPost, "https://api.sms.ir/v1/send/verify", bytes.NewBuffer(jsonSMS))
	if err != nil {
		return err
	}
	request.Header.Set("x-api-key", os.Getenv("SMS_API_KEY"))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}

	var smsOutput SMSOutput

	json.NewDecoder(response.Body).Decode(&smsOutput)

	if smsOutput.Status != 1 {
		return errors.New("SMS send error: " + smsOutput.Message)
	}

	return nil
}
