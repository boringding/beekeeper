//Mail interface.

package beekeeper

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"mime/quotedprintable"
	"net/mail"
	"net/smtp"
	"path/filepath"
	"strings"

	"github.com/boringding/beekeeper/mailmsg"
)

const (
	AddrTo = iota
	AddrCc
	AddrBcc
)

const (
	SmtpPort = 25
)

type loginAuth struct {
	userName string
	pwd      string
}

type Mail struct {
	from string
	to   []string
	msg  mailmsg.Msg
}

func (self *loginAuth) Start(serverInfo *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(self.userName), nil
}

func (self *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(self.userName), nil
		case "Password:":
			return []byte(self.pwd), nil
		}
	}
	return nil, nil
}

func (self *Mail) SetFrom(addr string) {
	self.from = addr
	self.msg.Header.From.Address = addr

	addrSlice := strings.Split(addr, "@")
	if len(addrSlice) > 1 {
		self.msg.Header.From.Name = addrSlice[0]
	}
}

func (self *Mail) AddReceiver(addr string, kind int) {
	self.to = append(self.to, addr)

	mailAddr := mail.Address{Address: addr}

	addrSlice := strings.Split(addr, "@")
	if len(addrSlice) > 1 {
		mailAddr.Name = addrSlice[0]
	}

	switch kind {
	case AddrTo:
		self.msg.Header.To = append(self.msg.Header.To, mailAddr)
	case AddrCc:
		self.msg.Header.Cc = append(self.msg.Header.Cc, mailAddr)
	}
}

func (self *Mail) SetSubject(subject string) {
	self.msg.Header.Subject = subject
}

func (self *Mail) SetContent(content string) error {
	self.msg.Body.CharSet = mailmsg.DefaultCharSet
	self.msg.Body.ContentType = mailmsg.TextPlain
	self.msg.Body.ContentTransferEncoding = mailmsg.EncodingBase64
	self.msg.Body.Content = content

	writer := &bytes.Buffer{}
	encoder := base64.NewEncoder(base64.StdEncoding, writer)
	_, err := encoder.Write([]byte(content))
	if err != nil {
		return err
	}

	err = encoder.Close()
	if err != nil {
		return err
	}

	self.msg.Body.Data = writer.Bytes()

	return nil
}

func (self *Mail) SetHtmlContent(content string) error {
	self.msg.HtmlBody.CharSet = mailmsg.DefaultCharSet
	self.msg.HtmlBody.ContentType = mailmsg.TextHtml
	self.msg.HtmlBody.ContentTransferEncoding = mailmsg.EncodingQuotedPrintable
	self.msg.HtmlBody.Content = content

	writer := &bytes.Buffer{}
	qWriter := quotedprintable.NewWriter(writer)

	_, err := qWriter.Write([]byte(content))
	if err != nil {
		return err
	}

	err = qWriter.Close()
	if err != nil {
		return err
	}

	self.msg.HtmlBody.Data = writer.Bytes()

	return nil
}

func (self *Mail) AddEmbeddedAttachment(path string) error {
	name := filepath.Base(path)

	embeddedAttachment := mailmsg.EmbeddedAttachment{Name: name, ContentType: mailmsg.ApplicationOctetStream, ContentId: name, ContentTransferEncoding: mailmsg.EncodingBase64}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	writer := &bytes.Buffer{}
	encoder := base64.NewEncoder(base64.StdEncoding, writer)
	_, err = encoder.Write(data)
	if err != nil {
		return err
	}

	err = encoder.Close()
	if err != nil {
		return err
	}

	embeddedAttachment.Data = writer.Bytes()

	self.msg.EmbeddedAttachments = append(self.msg.EmbeddedAttachments, embeddedAttachment)
	return nil
}

func (self *Mail) AddAttachment(path string) error {
	name := filepath.Base(path)

	attachment := mailmsg.Attachment{Name: name, ContentType: mailmsg.ApplicationOctetStream, FileName: name, ContentTransferEncoding: mailmsg.EncodingBase64}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	writer := &bytes.Buffer{}
	encoder := base64.NewEncoder(base64.StdEncoding, writer)
	_, err = encoder.Write(data)
	if err != nil {
		return err
	}

	err = encoder.Close()
	if err != nil {
		return err
	}

	attachment.Data = writer.Bytes()

	self.msg.Attachments = append(self.msg.Attachments, attachment)
	return nil
}

func LoginAuth(userName string, pwd string) smtp.Auth {
	return &loginAuth{userName: userName, pwd: pwd}
}

func NewMail() *Mail {
	return &Mail{}
}

func SendMail(host string, userName string, pwd string, mail *Mail) error {
	auth := LoginAuth(userName, pwd)
	msg := make([]byte, 0, 10000)
	var err error

	if (mail.msg.Body.Content != "" && mail.msg.HtmlBody.Content != "") || len(mail.msg.EmbeddedAttachments) > 0 || len(mail.msg.Attachments) > 0 {
		msg, err = mail.msg.WriteMime()
		if err != nil {
			return err
		}
	} else {
		msg, err = mail.msg.Write()
		if err != nil {
			return err
		}
	}

	addr := fmt.Sprintf("%s:%d", host, SmtpPort)

	err = smtp.SendMail(addr, auth, mail.from, mail.to, msg)
	if err != nil {
		return err
	}

	return nil
}
