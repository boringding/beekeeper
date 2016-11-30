package mailmsg

import (
	"fmt"
	"net/mail"
	"testing"
)

func Test_WriteMimeMsg(t *testing.T) {
	var em Msg

	em.Header.From = mail.Address{Name: "abc", Address: "abc@126.com"}
	em.Header.To = append(em.Header.To, mail.Address{Name: "efg", Address: "efg@126.com"})
	em.Header.To = append(em.Header.To, mail.Address{Name: "eee", Address: "eee@126.com"})
	em.Header.Cc = append(em.Header.Cc, mail.Address{Name: "hij", Address: "hij@126.com"})
	em.Header.Subject = "测试一下"

	em.Body.CharSet = DefaultCharSet
	em.Body.ContentType = "text/plain"
	em.Body.Data = []byte("这是一封测试邮件")
	em.Body.ContentTransferEncoding = "base64"

	em.HtmlBody.CharSet = DefaultCharSet
	em.HtmlBody.ContentType = "text/html"
	em.HtmlBody.Data = []byte("<h1>ADSFADF</h1>")
	em.HtmlBody.ContentTransferEncoding = "base64"

	em.EmbeddedAttachments = append(em.EmbeddedAttachments, EmbeddedAttachment{Name: "1.jpg", ContentId: "1.jpg", ContentType: ApplicationOctetStream, ContentTransferEncoding: "base64", Data: []byte("afasdfadfasdfasdfasdfasd")})
	em.EmbeddedAttachments = append(em.EmbeddedAttachments, EmbeddedAttachment{Name: "2.jpg", ContentId: "2.jpg", ContentType: ApplicationOctetStream, ContentTransferEncoding: "base64", Data: []byte("afasdfadfasdfasdfasdfasd")})

	em.Attachments = append(em.Attachments, Attachment{Name: "3.txt", ContentType: ApplicationOctetStream, ContentTransferEncoding: "base64", FileName: "3.txt"})
	em.Attachments = append(em.Attachments, Attachment{Name: "文件4.txt", ContentType: ApplicationOctetStream, ContentTransferEncoding: "base64", FileName: "文件4.txt"})

	msg, err := em.WriteMime()

	fmt.Println(string(msg))
	fmt.Println(err)
}
