package beekeeper

import (
	"fmt"
	"testing"
)

func Test_SendMimeMail(t *testing.T) {
	mail := NewMail()

	mail.SetFrom("example@example.com")
	mail.AddReceiver("dingyao@example.com", AddrTo)
	//mail.AddReceiver("vackyzhang@example.com", AddrCc)
	//mail.AddReceiver("dingyao@example.com", AddrBcc)

	mail.SetSubject("test")
	mail.SetContent("hahahahahahahaha")
	mail.SetHtmlContent(`<h1>fffff</h1><img src="cid:nav_logo242.png" />`)
	//mail.SetHtmlContent(`<h1>fffff</h1><h2>gggggg</h2>`)
	mail.AddEmbeddedAttachment(`C:\nav_logo242.png`)
	mail.AddAttachment(`C:\nav_logo242.png`)

	err := SendMail("smtp.example.com", "example", "xxxxxx", mail)
	if err != nil {
		fmt.Println(err)
	}
}
