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
	ID       string `json:"id"`
	From     string `json:"from"`
	To       string `json:"to"`
	Subject  string `json:"subject"`
	BodyHTML string `json:"body_html"`
	buffer   *bytes.Buffer
}

func (m *Message) setHeader(key, value string) {
	m.buffer.WriteString(key)
	m.buffer.WriteString(": ")
	m.buffer.WriteString(value)
	m.buffer.WriteString("\r\n")
}

func (m *Message) addBody(content string, contentType string) {
	m.buffer.WriteString("\r\n")
	m.buffer.WriteString(content)
	m.buffer.WriteString("\r\n")
	m.buffer.WriteString("\r\n")
}
