package config

type Config struct {
	AppName string
	Port int
	DBURI string
	JWTKEY string
	URL string
	Email EmailConfig 
	Phone PhoneConfig
	Oauth OAuthConfig
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

type OAuthConfig struct {
	GoogleClientID string
	GoogleSecretClientID string
	GithubClientID string
	GithubSecretClientId string
}

var AppConfig = &Config{
	AppName: "Snap_Backend_Practice",
	Port: 6060,
	DBURI: "mongodb+srv://abdrahman:abdrahman@rahmann18.hy9zl.mongodb.net/Snap_Backend_Practice",
	URL: "http://localhost:6060",
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
	Oauth: OAuthConfig{
		GoogleClientID: "",
		GoogleSecretClientID: "",
		GithubClientID: "",
		GithubSecretClientId: "",
	},
}