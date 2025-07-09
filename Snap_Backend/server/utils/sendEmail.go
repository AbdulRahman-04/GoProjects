package utils

import (
	"fmt"

	"github.com/AbdulRahman-04/GoProjects/Snap_Backend/config"
	"gopkg.in/gomail.v2"
)

type EmailData struct{
	From string
	To string
	Subject string
	Text string
	Html string
}

func SendEmail(data EmailData) error {
	// get user nd pass
	user := config.AppConfig.Email.User
	pass := config.AppConfig.Email.Pass

	// create sender 
	s := gomail.NewMessage()

	s.SetAddressHeader("From", user, "Team Snap Prac")
	s.SetHeader("To", data.To)
	s.SetHeader("Subject", data.Subject)
	s.SetBody("Text/plain", data.Text)
	s.AddAlternative("text/html", data.Html)

	// create transporter
	t := gomail.NewDialer("smtp.gmail.com", 465, user, pass)

	// try sending mail
	if err := t.DialAndSend(s); err != nil {
		fmt.Println("Error while sending email", err)
	}

	fmt.Println("Email Sent✅")
	return nil
}