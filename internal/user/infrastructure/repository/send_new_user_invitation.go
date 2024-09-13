package repository

import (
	"log"
	"os"

	externalUserDomain "github.com/goplaceapp/goplace-user/pkg/userservice/domain"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func (r *UserRepository) SendNewUserInvitationEmail(user *externalUserDomain.User, generatedPassword string) error {
	from := mail.NewEmail("GoPlace", os.Getenv("SENDGRID_NO_REPLY_EMAIL"))
	subject := "You're invited to join GoPlace!"
	to := mail.NewEmail(user.FirstName+" "+user.LastName, user.Email)

	plainTextContent := "Hello " + user.FirstName + " " + user.LastName + ",\n\n" +
		"You have been invited to join GoPlace. Please use the following credentials to login:\n\n" +
		"Your Email: " + user.Email + "\n" +
		"Your Password: " + generatedPassword + "\n\n" +
		"Your Pin Code: " + user.PinCode + "\n\n" +
		"Please click here to login: https://goplace.io/login\n\n" +
		"Thank you,\n" +
		"GoPlace Team"

	htmlContent := "Hello " + user.FirstName + " " + user.LastName + ",<br><br>" +
		"You have been invited to join GoPlace. Please use the following credentials to login:<br><br>" +
		"<strong>Your Email: " + user.Email + "</strong><br>" +
		"<strong>Your Password: " + generatedPassword + "</strong><br><br>" +
		"<strong>Your Pin Code: " + user.PinCode + "</strong>" +
		"<br><br>" +
		"Please click <a href='https://goplace.io/login'>here</a> to login." +
		"<br><br>" +
		"Thank you,<br>" +
		"GoPlace Team"

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	_, err := client.Send(message)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
