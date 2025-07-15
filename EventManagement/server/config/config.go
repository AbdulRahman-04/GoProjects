package config

type Config struct {
	AppName string
	Port    int
	DBURI   string
	URL     string
	JWTKEY  string
	Email   EmailConfig
	Phone   PhoneConfig
}

type EmailConfig struct {
	User string
	Pass string
}

type PhoneConfig struct {
	Sid   string
	Token string
	Phone string
}

var AppConfig = &Config{
	AppName: "Event_Booking",
	Port:    4040,
	DBURI:   "mongodb+srv://<username>:<password>@cluster.mongodb.net/Event_Booking",
	URL:     "http://localhost:4040",
	JWTKEY:  "your_jwt_secret_here",
	Email: EmailConfig{
		User: "your_email@example.com",
		Pass: "your_app_password_here",
	},
	Phone: PhoneConfig{
		Sid:   "your_twilio_sid_here",
		Token: "your_twilio_token_here",
		Phone: "+1234567890",
	},
}
