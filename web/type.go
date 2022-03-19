package web

type User struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type UserInfo struct {
	Account   string `json:"account"`
	Password  string `json:"password"`
	Signature string `json:"signature"`
}

type Task struct {
	ReceiverName string            `json:"receiverName"`
	EmailAddress string            `json:"emailAddress"`
	Location     TaskChildLocation `json:"location"`
	Time         TaskChildTime     `json:"time"`
	Weather      bool              `json:"weather"`
	News         bool              `json:"news"`
	Covid        bool              `json:"covid"`
}

type TaskChildLocation struct {
	Country  string `json:"country"`
	Province string `json:"province"`
	City     string `json:"city"`
}

type TaskChildTime struct {
	Hour   string `json:"hour"`
	Minute string `json:"minute"`
}

type Content struct {
	ReceiverName string `json:"receiverName"`
	EmailAddress string `json:"emailAddress"`
	Content      string `json:"content"`
}
