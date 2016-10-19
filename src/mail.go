package main

import (
	"net/smtp"
	"fmt"
)


func SendEmail(to string, subject string, body string) {
	config, _ := Configuration()
  from := config.Email
	pass := config.Password
	cc := config.Email
	fmt.Println("Report User or Comment to " + from)
	msg := "From: " + from + "\n" +
		"To: " + from + "\n" +
    "Cc: " + cc + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))
}
