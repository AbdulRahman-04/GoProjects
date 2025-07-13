package utils

import (
	"fmt"

	"github.com/AbdulRahman-04/GoProjects/RestaurantManagement/server/config"
	"gopkg.in/gomail.v2"
)

type EmailData struct {
	From string
	To string
	Subject string
	Text string
	Html string
}

func SendEmail(data EmailData) error {
	// get user and pass from config
	user := config.AppConfig.Email.User
	pass := config.AppConfig.Email.Pass

	// create new message 
	s := gomail.NewMessage()

	s.SetAddressHeader("From", user, "Team Restro Management")
	s.SetBody("To", data.To)
	s.SetBody("From", data.From)
	s.SetBody("Text/plain", data.Text)
	s.AddAlternative("text/html", data.Html)

	// create smtp 
	t := gomail.NewDialer("smtp.gmail.com", 465, user, pass)

	// try sending mail 
	if err := t.DialAndSend(s); err != nil {
		fmt.Println(err)
	}

	fmt.Println("Email Sent✅")
	return  nil
}