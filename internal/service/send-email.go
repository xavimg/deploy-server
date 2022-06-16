package service

import (
	"log"
	"math/rand"
	"net/smtp"
	"strconv"
	"time"
)

const (
	email    string = "alanturingoffworld@gmail.com"
	password string = "alanturing123456"
)

func SendEmail(username string, To string) {

	from := email
	password := password
	toEmail := To
	to := []string{toEmail}

	host := "smtp.gmail.com"
	port := "587"
	address := host + ":" + port

	auth := smtp.PlainAuth("", from, password, host)
	msg := []byte(
		"From: Off World <" + from + ">\r\n" +
			"To: " + toEmail + "\r\n" +
			"Subject: Off World Welcoming!\r\n" +
			"MIME: MIME-version: 1.0\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\";\r\n" +
			"\r\n" +
			`<html>
				<h1>Welcome to <b>Alan Turing</b> family</h1>

				<p> If you've received this mail, it means that you are a true Space explorer <p>

				<p> You are pre-registered for our next game Off World <p>

				<p> The downloable beta will be realesed at <b> 01-05-2002 <b> <p>

				<h3> Team: <h3>
				<ul>
					<li>Khadija Rehman</li>
					<li>Alex Andreba</li>
					<li>Gerard Marquina</li>
					<li>Xavier Moya</li>
				</ul>
			</html>`)

	err := smtp.SendMail(address, auth, from, to, msg)
	if err != nil {
		log.Fatal(err)
	}

}

func SendEmailCodeVerify(username string, To string) int {

	from := email
	password := password
	toEmail := To
	to := []string{toEmail}

	host := "smtp.gmail.com"
	port := "587"
	address := host + ":" + port

	codeVerify := GenerateVerificationCode()

	auth := smtp.PlainAuth("", from, password, host)

	message := ([]byte("This is your code verification: " + strconv.Itoa(codeVerify)))

	msg := []byte(
		"From: Off World <" + from + ">\r\n" +
			"To: " + toEmail + "\r\n" +
			"Subject: Verify this code for start playing!\r\n" +
			"MIME: MIME-version: 1.0\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\";\r\n" +
			"\r\n" +
			string(message))

	err := smtp.SendMail(address, auth, from, to, msg)
	if err != nil {
		log.Fatal(err)
	}

	return codeVerify
}

func GenerateVerificationCode() (code int) {

	max := 999999
	min := 100000

	rand.Seed(time.Now().Unix())

	return rand.Intn(max-min) + min

}
