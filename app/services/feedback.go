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
	message      structs.Message
	captchaKey   string
	smtpServer   string
}

func NewFeedbackService(cfg config.SMTP) (FeedbackService, error) {
	auth := smtp.PlainAuth(
		"",
		cfg.MailAccount,
		cfg.AccountPassword,
		cfg.SMTPServerAddress,
	)

	template, err := template.ParseFiles(templatPath)
	if err != nil {
		return FeedbackService{}, fmt.Errorf("error in PasreFiles(): %w", err)
	}

	message := structs.Message{
		From:    cfg.MailAccount,
		To:      cfg.MailAccount,
		Subject: "AtakHelp Feedback",
	}

	return FeedbackService{
		auth:         auth,
		templateFile: template,
		message:      message,
		captchaKey:   cfg.CaptchaKey,
		smtpServer:   cfg.SMTPServerAddress,
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

	var body bytes.Buffer
	if err := f.templateFile.Execute(&body, feedback); err != nil {
		return fmt.Errorf("error in templateFile.Execute(): %w", err)
	}

	if err := f.sendMail(body); err != nil {
		return err
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

func (f FeedbackService) sendMail(body bytes.Buffer) error {
	f.message.Buffer = bytes.NewBuffer(make([]byte, 256))
	f.message.Buffer.Reset()
	f.message.SetHeader("MIME-Version", "1.0")

	f.message.SetHeader("From", f.message.From)
	f.message.SetHeader("To", f.message.To)
	f.message.SetHeader("Subject", f.message.Subject)
	f.message.SetHeader("Content-Type", "text/html; charset=UTF-8")
	f.message.Buffer.WriteString("\r\n")
	f.message.AddBody(body.String(), "text/html")

	err := smtp.SendMail(
		f.smtpServer+":587",
		f.auth,
		f.message.From,
		[]string{f.message.To},
		f.message.Buffer.Bytes(),
	)
	if err != nil {
		return fmt.Errorf("error in smtp.SendMai(): %w", err)
	}

	return nil
}
