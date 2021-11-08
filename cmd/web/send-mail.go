package main

import (
	"fmt"
	mail "github.com/xhit/go-simple-mail/v2"
	"io/ioutil"
	"log"
	"myapp/internal/models"
	"strings"
	"time"
)

func listenForMail() {
	go func() {
		for {
			msg := <- app.MailChan
			sendMsg(msg)
		}
	}()
}

func sendMsg(m models.MailData) {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025 // real mail servers use ports: 25 || 587 || 465
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	client, err := server.Connect()
	if err != nil {
		errorLog.Println(err)
	}

	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	if m.Template == "" {
		email.SetBody(mail.TextHTML, m.Content)
	} else {
		data, err := ioutil.ReadFile(fmt.Sprintf("./email-templates/%s", m.Template))
		if err != nil {
			app.ErrorLog.Println(err)
		}

		mailTemplate := string(data)
		msgToSend := strings.Replace(mailTemplate, "[%body%]", m.Content, 1)
		email.SetBody(mail.TextHTML, msgToSend)
	}

	err = email.Send(client) // client contains the email
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Email sent!")
	}
}