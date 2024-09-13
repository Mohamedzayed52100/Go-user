package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/goplaceapp/goplace-common/pkg/logger"
	"github.com/goplaceapp/goplace-user/pkg/userservice/domain"
)

func (r *UserRepository) SendResetPasswordWhatsappMessage(ctx context.Context, currentUser *domain.User, otp string) error {
	requestBody := map[string]interface{}{
		"channelId":   os.Getenv("GALLABOX_CHANNEL_ID"),
		"channelType": "whatsapp",
		"recipient": map[string]string{
			"name":  currentUser.FirstName + " " + currentUser.LastName,
			"phone": currentUser.PhoneNumber,
		},
		"whatsapp": map[string]interface{}{
			"type": "template",
			"template": map[string]interface{}{
				"templateName": "reset_password_v2",
				"bodyValues": map[string]string{
					"name": currentUser.FirstName,
					"otp":  otp,
				},
			},
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://server.gallabox.com/devapi/messages/whatsapp", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apiKey", os.Getenv("GALLABOX_API_KEY"))
	req.Header.Set("apiSecret", os.Getenv("GALLABOX_API_SECRET"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		logger.Default().Error("Failed to send whatsapp message to: ", requestBody["recipient"], " with status code: ", resp.StatusCode, " and response: ", resp.Body)
	}

	return nil
}
