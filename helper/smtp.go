package helper

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
)

type SmtpHelper struct {
	host     string
	port     string
	username string
	password string
	identity string
}

type SmtpHelperOpts struct {
	Host     string
	Port     string
	Username string
	Password string
	Identity string
}

func NewSmptpHelper(opt SmtpHelperOpts) *SmtpHelper {
	return &SmtpHelper{
		host:     opt.Host,
		port:     opt.Port,
		username: opt.Username,
		password: opt.Password,
		identity: opt.Identity,
	}
}

type smtpRequest struct {
	helper *SmtpHelper

	username   string
	recipients []string
	subject    string
	body       string
}

func (h *SmtpHelper) NewRequest(to []string, subject string) *smtpRequest {
	return &smtpRequest{
		helper: h,

		username:   h.username,
		recipients: to,
		subject:    subject,
	}
}

func (r *smtpRequest) SetBody(emailTemplate string, data any) (*smtpRequest, error) {
	t, err := template.New("emailTemplate").Parse(emailTemplate)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)

	err = t.Execute(buf, data)
	if err != nil {
		return nil, err
	}

	r.body = buf.String()

	return r, nil
}

func (r *smtpRequest) Send() error {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + r.subject + "!\n"
	msg := []byte(subject + mime + "\n" + r.body)

	addr := fmt.Sprintf("%s:%s", r.helper.host, r.helper.port)

	auth := smtp.PlainAuth(r.helper.identity, r.helper.username, r.helper.password, r.helper.host)

	if err := smtp.SendMail(addr, auth, r.username, r.recipients, msg); err != nil {
		return err
	}

	return nil
}
