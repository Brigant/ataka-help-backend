package services

import (
	"bytes"
	"context"
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
	auth         smtp.Auth
	templateFile *template.Template
	captchaKey   string
}

func NewFeedbackService(cfg config.SMTP) (FeedbackService, error) {
	auth := smtp.PlainAuth("", cfg.MailAccount, cfg.AccountPassword, cfg.SMTPServerAddress)

	template, err := template.ParseFiles(templatPath)
	if err != nil {
		return FeedbackService{}, fmt.Errorf("error in checkGoogleCaptcha(): %w", err)
	}

	return FeedbackService{
		auth:         auth,
		templateFile: template,
		captchaKey:   cfg.CaptchaKey,
	}, nil
}

func (f FeedbackService) PassFeedback(ctx context.Context, feedback structs.Feedback) error {
	ok, err := f.checkGoogleCaptcha(ctx, feedback.Token, f.captchaKey)
	if err != nil {
		return fmt.Errorf("error in checkGoogleCaptcha(): %w", err)
	}

	if !ok {
		return structs.ErrCheckCaptcha
	}

	return nil
}

func (f FeedbackService) checkGoogleCaptcha(ctx context.Context, token, googleCaptcha string) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://www.google.com/recaptcha/api/siteverify", nil)
	if err != nil {
		return false, fmt.Errorf("error in http.NewRequest(): %w", err)
	}

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

	isValid, ok := googleResponse["success"].(bool)
	if !ok {
		return false, fmt.Errorf("error in the assertion: %w", err)
	}

	return isValid, nil
}

func (f FeedbackService) sMail() {
}

func (f FeedbackService) sendMail(feedback structs.Feedback, mailAccount, mailPassword, smtpServer string, template *template.Template) error {
	var body bytes.Buffer

	if err := template.Execute(&body, feedback); err != nil {
		return fmt.Errorf("error in templateFile.Execute(): %w", err)
	}

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

	err := smtp.SendMail(
		smtpServer+":587",
		auth,
		mailAccount,
		[]string{mailAccount},
		[]byte(msg),
	)
	if err != nil {
		return fmt.Errorf("error in smtp.SendMai(): %w", err)
	}

	return nil
}
