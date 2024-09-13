package repository

import (
	"os"

	"github.com/goplaceapp/goplace-user/pkg/userservice/domain"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func (r *UserRepository) SendResetPasswordEmail(user *domain.User, token, otp string) (bool, error) {
	from := mail.NewEmail("GoPlace", os.Getenv("SENDGRID_NO_REPLY_EMAIL"))
	subject := "Reset Password - GoPlace"
	to := mail.NewEmail(user.FirstName+" "+user.LastName, user.Email)

	pinCode := otp
	pinCodePageLink := "https://goplace.io/reset-password/" + token + "/verify?otp=" + pinCode

	plainTextContent := "Hello " + user.FirstName + " " + user.LastName + ",\n\n" +
		"You have requested to reset your password. Please use the following pin code to reset your password:\n\n" +
		pinCode + "\n\n" +
		"Or you can directly go to the reset password page by clicking the link below:\n\n" +
		pinCodePageLink + "\n\n" +
		"Thank you,\n" +
		"GoPlace Team"

	htmlContent := "Hello " + user.FirstName + " " + user.LastName + ",<br><br>" +
		"You have requested to reset your password. Please use the following pin code to reset your password:<br><br>" +
		"<b>" + pinCode + "</b><br><br>" +
		"Or you can directly go to the reset password page by clicking the link below:<br><br>" +
		"<a href='" + pinCodePageLink + "'>Reset Password</a><br><br>" +
		"Thank you,<br>" +
		"GoPlace Team"

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	_, err := client.Send(message)
	if err != nil {
		return false, err
	}

	return true, nil
}
