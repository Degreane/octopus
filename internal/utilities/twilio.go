package utilities

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
	lua "github.com/yuin/gopher-lua"
)

// Configuration struct to hold Twilio credentials from environment variables
type TwilioConfig struct {
	AccountSID string
	AuthToken  string
	FromNumber string
}

// LoadTwilioConfig loads Twilio configuration from environment variables
func LoadTwilioConfig() (*TwilioConfig, error) {
	accountSID := os.Getenv("TWILIO_ACCOUNT_SID")

	authToken := os.Getenv("TWILIO_AUTH_TOKEN")

	fromNumber := os.Getenv("TWILIO_FROM_NUMBER") // Add from number

	if accountSID == "" || authToken == "" || fromNumber == "" {
		return nil, fmt.Errorf("missing required Twilio environment variables")
	}

	return &TwilioConfig{
		AccountSID: accountSID,
		AuthToken:  authToken,
		FromNumber: fromNumber,
	}, nil
}

// SendWhatsAppMessage sends a WhatsApp message using the provided message string and recipient number.
// It uses the Twilio credentials from the TwilioConfig.
func SendWhatsAppMessage(toNumber, message string) error {
	config, err := LoadTwilioConfig()
	if err != nil {
		return err
	}
	// Validate the "To" number format for WhatsApp
	if !isValidWhatsAppNumber(toNumber) {
		return fmt.Errorf("invalid WhatsApp 'To' number format. Must be 'whatsapp:+1XXXXXXXXXX'")
	}

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: config.AccountSID,
		Password: config.AuthToken,
	})

	//Format the number for whatsapp sending
	fromNumber := fmt.Sprintf("whatsapp:%s", config.FromNumber)
	toNumber = fmt.Sprintf("whatsapp:%s", toNumber)
	params := &api.CreateMessageParams{}
	// params.SetMessagingServiceSid("HX5f78dc700628e01aec1970ce264506d5")
	params.SetMessagingServiceSid("MG20d448746827f963a5b9e32c243f5b10")
	params.SetContentSid("HXb00a2452de987a7934b9e2bd0fa72bfe")
	cv, err := json.Marshal(map[string]interface{}{
		"1": message,
	})
	params.SetContentVariables(string(cv))
	params.SetBody(message)
	params.SetFrom(fromNumber)
	params.SetTo(toNumber)

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("error sending WhatsApp message: %w", err)
	}

	if err != nil { // Explicitly check the status code
		return fmt.Errorf("failed to send message Twilio API returned status: %s, body: %s", resp.Status, *resp.Body)
	}
	response, _ := json.Marshal(*resp)
	log.Printf("Twilio response: %s", string(response))
	// log.Printf("Twilio response: %+v", resp.Body)

	return nil
}

// isValidWhatsAppNumber checks if the given phone number string is in the correct WhatsApp format.
func isValidWhatsAppNumber(number string) bool {
	// Add your WhatsApp number validation logic here.
	// This is a placeholder. You might want to use a regular expression or more robust validation.
	// For this example, it just checks if the string starts with "whatsapp:+"
	return regexp.MustCompile(`^\+[1-9]\d{1,}$`).MatchString(number)
}

// SendWhatsAppMessageLua is a wrapper function for gopher-lua binding
func SendWhatsAppMessageLua(L *lua.LState) int {
	toNumber := L.CheckString(1)
	message := L.CheckString(2)

	err := SendWhatsAppMessage(toNumber, message)
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}

	return 0
}
