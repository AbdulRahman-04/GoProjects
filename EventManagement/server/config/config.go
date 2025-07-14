package config

type Config struct {
	AppName string
	Port int
	DBURI string
	URL string
	JWTKEY string
	Email EmailConfig 
	Phone PhoneConfig
}

type EmailConfig struct {
	User string
	Pass string
}

type PhoneConfig struct {
	Sid string
	Token string
	Phone string
}

var AppConfig = &Config{
	AppName: "Event_Booking",
	Port: 4040,
	DBURI: "mongodb+srv://abdrahman:abdrahman@rahmann18.hy9zl.mongodb.net/Event_Booking",
	URL: "http://localhost:4040",
	JWTKEY: "RAHMAN123",
	Email: EmailConfig{
		User: "abdulrahman.81869@gmail.com",
		Pass: "znoh cwef huhl hvln",
	},
	Phone: PhoneConfig{
		Sid: "",
		Token: "",
		Phone: "",
	},
}