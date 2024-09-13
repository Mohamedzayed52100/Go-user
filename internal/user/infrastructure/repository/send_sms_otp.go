package repository

import (
	"fmt"
	"os"

	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"

	"github.com/goplaceapp/goplace-user/pkg/userservice/domain"
	"github.com/twilio/twilio-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (r *UserRepository) SendSMSOTP(getUser *domain.User, otp string) error {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: os.Getenv("TWILIO_ACCOUNT_SID"),
		Password: os.Getenv("TWILIO_AUTH_TOKEN"),
	})

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(getUser.PhoneNumber)
	params.SetFrom(os.Getenv("TWILIO_PHONE_NUMBER"))
	params.SetBody(fmt.Sprintf("Hi, %s\nYour OTP is %s", getUser.FirstName, otp))

	if _, err := client.Api.CreateMessage(params); err != nil {
		return status.Error(codes.Internal, "Error sending otp: "+err.Error())
	}

	return nil
}
