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
	AppName: "Restro_Management",
	Port: 5080,
	DBURI: "mongodb+srv://abdrahman:abdrahman@rahmann18.hy9zl.mongodb.net/Restro_Management",
	URL: "http://localhost:5080",
	JWTKEY: "RAHMAN123",
	Email: EmailConfig{
		User: "abdulrahman.81869@gmail.com",
		Pass: "mliw uitm xunh usjr",
	},
	Phone: PhoneConfig{
		Sid: "",
		Token: "",
		Phone: "",
	},
}