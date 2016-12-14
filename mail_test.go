package beekeeper

import (
	"fmt"
	"testing"
)

func Test_SendMimeMail(t *testing.T) {
	mail := NewMail()

	mail.SetFrom("example@example.com")
	mail.AddReceiver("example@example.com", AddrTo)
	//mail.AddReceiver("example@example.com", AddrCc)
	//mail.AddReceiver("example@example.com", AddrBcc)

	mail.SetSubject("test")
	mail.SetContent("hahahahahahahaha")
	mail.SetHtmlContent(`<h1>fffff</h1><img src="cid:googlelogo_color_120x44dp.png" />`)
	//mail.SetHtmlContent(`<h1>fffff</h1><h2>gggggg</h2>`)
	mail.AddEmbeddedAttachment(`C:\googlelogo_color_120x44dp.png`)
	mail.AddAttachment(`C:\googlelogo_color_120x44dp.png`)

	err := SendMail("smtp.example.com", "example", "example", mail)
	if err != nil {
		fmt.Println(err)
	}
}
