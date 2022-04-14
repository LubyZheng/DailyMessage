package email

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/mail"
	"net/smtp"
	"os"
	"time"
)

const EMAIL_FORMAT = `
Dear {{ReceiverName}}, 
<br>
<br>
{{Content}} 
<br>
<br>
Best regards,
<br>
<br>
{{Signature}}
`

func parseMailAddr(address string) *mail.Address {
	addr, err := mail.ParseAddress(address)
	if err != nil {
		log.Fatalf("Parse email address %s error: %s", address, err)
	}
	return addr
}

func (m *Mail) String() string {
	var buf bytes.Buffer
	const crlf = "\r\n"

	write := func(what string, addrs []string) {
		if len(addrs) == 0 {
			return
		}
		for i := range addrs {
			if i == 0 {
				buf.WriteString(what)
			} else {
				buf.WriteString(", ")
			}
			buf.WriteString(parseMailAddr(addrs[i]).String())
		}
		buf.WriteString(crlf)
	}
	getBoundary := func() string {
		h := md5.New()
		io.WriteString(h, fmt.Sprintf("%d", time.Now().Nanosecond()))
		return fmt.Sprintf("%x", h.Sum(nil))
	}

	from := parseMailAddr(m.From)
	if from.Address == "" {
		from = parseMailAddr(Config.From)
	}
	fmt.Fprintf(&buf, "From: %s%s", from.String(), crlf)
	write("To: ", m.To)
	write("Cc: ", m.Cc)
	write("Bcc: ", m.Bcc)
	boundary := getBoundary()
	fmt.Fprintf(&buf, "Date: %s%s", time.Now().UTC().Format(time.RFC822), crlf)
	fmt.Fprintf(&buf, "Subject: %s%s", m.Subject, crlf)
	fmt.Fprintf(&buf, "Content-Type: multipart/alternative; boundary=%s%s%s", boundary, crlf, crlf)
	fmt.Fprintf(&buf, "%s%s", "--"+boundary, crlf)
	fmt.Fprintf(&buf, "Content-Type: text/HTML; charset=UTF-8%s", crlf)
	fmt.Fprintf(&buf, "%s%s%s%s", crlf, m.Content, crlf, crlf)
	fmt.Fprintf(&buf, "%s%s", "--"+boundary+"--", crlf)

	return buf.String()
}

// Send email
func (m *Mail) Send() error {
	to := make([]string, len(m.To))
	for i := range m.To {
		to[i] = parseMailAddr(m.To[i]).Address
	}

	if m.From == "" {
		m.From = Config.From
	}
	from := parseMailAddr(m.From).Address
	addr := fmt.Sprintf("%s:%s", Config.Host, Config.Port)
	auth := smtp.PlainAuth("", Config.Username, Config.Password, Config.Host)
	return smtp.SendMail(addr, auth, from, to, []byte(m.String()))
}

func SendMail(content string, to string) {
	Config.Username = os.Getenv("MAIL_USERNAME")
	Config.Password = os.Getenv("MAIL_PASSWORD")
	Config.Host = os.Getenv("MAIL_HOST")
	Config.Port = os.Getenv("MAIL_PORT")
	Config.From = os.Getenv("MAIL_SIGNATURE") + "<" + os.Getenv("MAIL_USERNAME") + ">"

	email := Mail{
		To:      []string{to},
		Subject: os.Getenv("MAIL_SUBJECT"),
		Content: content,
	}

	err := email.Send()
	if err != nil {
		log.Printf("Send email fail, error: %s", err)
	} else {
		log.Printf("Send email %s success!", to)
	}
}
