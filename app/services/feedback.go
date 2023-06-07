package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/smtp"
	"text/template"

	"github.com/baza-trainee/ataka-help-backend/app/config"
	"github.com/baza-trainee/ataka-help-backend/app/structs"
)

type FeedbackService struct {
	cfg config.Smtp
}

func (f FeedbackService) PassFeedback(feedback structs.Feedback) error {
	ok, err := f.checkGoogleCaptcha(feedback.Token, f.cfg.CaptchaKey)
	if err != nil {
		return fmt.Errorf("error in checkGoogleCaptcha(): %w", err)
	}

	if !ok {
		return structs.ErrCheckCaptcha
	}

	if err := f.sendMail(feedback, f.cfg.MailAccount, f.cfg.AccountPassword, f.cfg.SmtpServerAddress); err != nil {
		return fmt.Errorf("error in sendMail(): %w", err)
	}

	return nil
}

func (f FeedbackService) checkGoogleCaptcha(token, googleCaptcha string) (bool, error) {
	req, _ := http.NewRequest("POST", "https://www.google.com/recaptcha/api/siteverify", nil)

	pathQuery := req.URL.Query()
	pathQuery.Add("secret", googleCaptcha)
	pathQuery.Add("response", token)

	req.URL.RawQuery = pathQuery.Encode()

	client := &http.Client{}

	var googleResponse map[string]interface{}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("error in io.ReadAll(): %w", err)
	}

	if err := json.Unmarshal(body, &googleResponse); err != nil {
		return false, fmt.Errorf("error in json.Unmarshal(): %w", err)
	}

	return googleResponse["success"].(bool), nil
}

func (f FeedbackService) sendMail(feedback structs.Feedback, mailAccount, mailPassword, smtpServer string) error {
	var body bytes.Buffer

	templateFile, err := template.ParseFiles("./static/template/emailTemplate.html")
	if err != nil {
		return fmt.Errorf("error in template.ParseFiles(): %w", err)
	}

	templateFile.Execute(&body, feedback)

	auth := smtp.PlainAuth(
		"",
		mailAccount,
		mailPassword,
		smtpServer,
	)

	headers := fmt.Sprintf(
		"MIME-version: 1.0;\n"+
			"Return-Path: <\"%s\">\n"+
			"From: \"%s\";\n"+
			"To: \"%s\";\n"+
			"Content-Type: text/html; charset=\"UTF-8\";",
		mailAccount,
		mailAccount,
		mailAccount,
	)

	msg := "Subject: AtackHelp Feedback\n" + headers + "\n\n" + body.String()

	err = smtp.SendMail(
		smtpServer+":587",
		auth,
		"sergey.july@nis.ua",
		[]string{"sergey.july@nis.ua"},
		[]byte(msg),
	)

	if err != nil {
		return fmt.Errorf("error in smtp.SendMai(): %w", err)
	}

	return nil
}
