package beekeeper

import (
	"fmt"
	"testing"
)

func Test_SendMimeMail(t *testing.T) {
	mail := NewMail()

	mail.SetFrom("wetest@tencent.com")
	mail.AddReceiver("dingyao@tencent.com", AddrTo)
	//mail.AddReceiver("vackyzhang@tencent.com", AddrCc)
	//mail.AddReceiver("dingyao@tencent.com", AddrBcc)

	mail.SetSubject("daxi")
	mail.SetContent("hahahahahahahaha")
	mail.SetHtmlContent(`<h1>fffff</h1><img src="cid:nav_logo242.png" />`)
	//mail.SetHtmlContent(`<h1>fffff</h1><h2>gggggg</h2>`)
	mail.AddEmbeddedAttachment(`C:\Users\dingyao\Desktop\nav_logo242.png`)
	mail.AddAttachment(`C:\Users\dingyao\Desktop\nav_logo242.png`)

	err := SendMail("smtp.tencent.com", "wetest", "WtServer@Ieg930", mail)
	if err != nil {
		fmt.Println(err)
	}
}
