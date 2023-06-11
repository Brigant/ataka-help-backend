package structs

import (
	"bytes"
	"errors"
	"regexp"
)

const (
	nameMask    = `^[\p{L}\s'’]{2,50}$`
	mailMask    = `^[\w-\.]+@[\w-]+\.+[\w-]{2,20}$`
	badMask     = `@[a-zA-Z0-9.-]*\.ru$`
	commentMask = `^[\p{L}\s'’\(\)?!.,\-]{1,300}$`
)

var (
	nameRegex       = regexp.MustCompile(nameMask)
	mailRegex       = regexp.MustCompile(mailMask)
	commnentRegex   = regexp.MustCompile(commentMask)
	exlcludeRegex   = regexp.MustCompile(badMask)
	errWrongName    = errors.New("wrong name")
	errWrongEmail   = errors.New("wrong email")
	errWrongComment = errors.New("wrong comment")
	errToken        = errors.New("empty toker")
)

type Feedback struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Comment string `json:"comment"`
	Token   string `json:"token"`
}

func (f Feedback) Valiadate() error {
	if !nameRegex.MatchString(f.Name) {
		return errWrongName
	}

	if !mailRegex.MatchString(f.Email) || exlcludeRegex.MatchString(f.Email) {
		return errWrongEmail
	}

	if !commnentRegex.MatchString(f.Comment) {
		return errWrongComment
	}

	if len(f.Token) < 1 {
		return errToken
	}

	return nil
}

type Message struct {
	From     string
	To       string
	Subject  string
	BodyHTML []byte
	Buffer   *bytes.Buffer
}

func (m *Message) SetHeader(key, value string) {
	m.Buffer.WriteString(key)
	m.Buffer.WriteString(": ")
	m.Buffer.WriteString(value)
	m.Buffer.WriteString("\r\n")
}

func (m *Message) AddBody(content string, contentType string) {
	m.Buffer.WriteString("\r\n")
	m.Buffer.WriteString(content)
	m.Buffer.WriteString("\r\n")
	m.Buffer.WriteString("\r\n")
}
