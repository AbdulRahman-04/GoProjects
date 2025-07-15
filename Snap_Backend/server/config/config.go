package config

type Config struct {
	AppName string
	Port    int
	DBURI   string
	JWTKEY  string
	URL     string
	Email   EmailConfig
	Phone   PhoneConfig
	Oauth   OAuthConfig
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

type OAuthConfig struct {
	GoogleClientID        string
	GoogleSecretClientID  string
	GithubClientID        string
	GithubSecretClientId  string
}

var AppConfig = &Config{
	AppName: "Snap_Backend_Practice",
	Port:    6060,
	DBURI:   "mongodb+srv://<username>:<password>@cluster.mongodb.net/Snap_Backend_Practice",
	URL:     "http://localhost:6060",
	JWTKEY:  "your_jwt_secret_key",
	Email: EmailConfig{
		User: "your_email@example.com",
		Pass: "your_email_app_password",
	},
	Phone: PhoneConfig{
		Sid:   "your_twilio_sid",
		Token: "your_twilio_token",
		Phone: "+1234567890",
	},
	Oauth: OAuthConfig{
		GoogleClientID:        "your_google_client_id",
		GoogleSecretClientID:  "your_google_secret",
		GithubClientID:        "your_github_client_id",
		GithubSecretClientId:  "your_github_secret",
	},
}
