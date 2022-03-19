package email

// Configuration for mail
type Configuration struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

// Config default configuration
var Config = Configuration{
	Host:     "smtp.qq.com",
	Port:     "25",
	Username: "",
	Password: "",
	From:     "",
}

// Mail config
type Mail struct {
	From    string
	To      []string
	Cc      []string
	Bcc     []string
	Subject string
	Content string
}
