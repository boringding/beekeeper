//Internal representation of smtp mail message.

package mailmsg

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/mail"
	"net/textproto"
)

type Header struct {
	From        mail.Address
	To          []mail.Address
	Cc          []mail.Address
	Subject     string
	ContentType string
}

type Body struct {
	CharSet                 string
	ContentType             string
	ContentTransferEncoding string
	Content                 string
	Data                    []byte
}

type EmbeddedAttachment struct {
	Name                    string
	ContentType             string
	ContentId               string
	ContentTransferEncoding string
	Data                    []byte
}

type Attachment struct {
	Name                    string
	ContentType             string
	FileName                string
	ContentTransferEncoding string
	Data                    []byte
}

type Msg struct {
	Header              Header
	Body                Body
	HtmlBody            Body
	EmbeddedAttachments []EmbeddedAttachment
	Attachments         []Attachment
}

const (
	DefaultCharSet          = "utf-8"
	DefaultMimeVersion      = "1.0"
	DefaultEncoding         = "base64"
	EncodingBase64          = "base64"
	EncodingQuotedPrintable = "quoted-printable"
)

const (
	MultipartMixed         = "multipart/mixed"
	MultipartRelated       = "multipart/related"
	MultipartAlternative   = "multipart/alternative"
	TextPlain              = "text/plain"
	TextHtml               = "text/html"
	ApplicationOctetStream = "application/octet-stream"
	DispositionAttachment  = "attachment"
)

func (self *Header) Write(writer io.Writer, boundary string) error {
	if len(self.To) <= 0 {
		return errors.New("empty to list")
	}

	if self.ContentType == "" {
		return errors.New("empty content type")
	}

	err := checkAddrs(self.From)
	err = checkAddrs(self.To...)
	err = checkAddrs(self.Cc...)
	if err != nil {
		return err
	}

	mimeHeader := textproto.MIMEHeader{}
	mimeHeader.Add("From", stringAddrs(self.From))
	mimeHeader.Add("To", stringAddrs(self.To...))

	if len(self.Cc) > 0 {
		mimeHeader.Add("Cc", stringAddrs(self.Cc...))
	}

	if self.Subject != "" {
		mimeHeader.Add("Subject", mime.BEncoding.Encode(DefaultCharSet, self.Subject))
	}

	mimeHeader.Add("Mime-Version", DefaultMimeVersion)
	mimeHeader.Add("Content-Type", fmt.Sprintf(`%s;%s boundary="%s"`, self.ContentType, Seperator, boundary))

	err = writeHeader(mimeHeader, writer)
	if err != nil {
		return err
	}

	return nil
}

func (self *Msg) Write() ([]byte, error) {
	writer := &bytes.Buffer{}

	if self.Body.Content != "" {
		self.Header.ContentType = TextPlain
	} else if self.HtmlBody.Content != "" {
		self.Header.ContentType = TextHtml
	} else {
		return writer.Bytes(), errors.New("empty body")
	}

	err := self.Header.Write(writer, "")
	if err != nil {
		return writer.Bytes(), err
	}

	if self.Body.Content != "" {
		_, err = writer.Write([]byte(self.Body.Content))
	} else if self.HtmlBody.Content != "" {
		_, err = writer.Write([]byte(self.HtmlBody.Content))
	}

	if err != nil {
		return writer.Bytes(), err
	}

	return writer.Bytes(), nil
}

func (self *Msg) WriteMime() ([]byte, error) {
	writer := &bytes.Buffer{}
	multipartWriter := multipart.NewWriter(writer)
	lastBytes := make([]byte, 0, 10000)
	lastBoundary := ""

	if len(self.Body.Data) > 0 {
		self.Header.ContentType = TextPlain

		mimeHeader := textproto.MIMEHeader{}

		mimeHeader.Add("Content-Type", fmt.Sprintf(`%s;%s charset="%s"`, self.Body.ContentType, Seperator, self.Body.CharSet))
		mimeHeader.Add("Content-Transfer-Encoding", self.Body.ContentTransferEncoding)

		w, err := multipartWriter.CreatePart(mimeHeader)
		if err != nil {
			return writer.Bytes(), err
		}

		_, err = w.Write(self.Body.Data)
		if err != nil {
			return writer.Bytes(), err
		}
	}

	if len(self.HtmlBody.Data) > 0 {
		self.Header.ContentType = TextHtml

		mimeHeader := textproto.MIMEHeader{}

		mimeHeader.Add("Content-Type", fmt.Sprintf(`%s;%s charset="%s"`, self.HtmlBody.ContentType, Seperator, self.HtmlBody.CharSet))
		mimeHeader.Add("Content-Transfer-Encoding", self.HtmlBody.ContentTransferEncoding)

		w, err := multipartWriter.CreatePart(mimeHeader)
		if err != nil {
			return writer.Bytes(), err
		}

		_, err = w.Write(self.HtmlBody.Data)
		if err != nil {
			return writer.Bytes(), err
		}
	}

	if len(self.Body.Data) > 0 && len(self.HtmlBody.Data) > 0 {
		self.Header.ContentType = MultipartAlternative

		if len(self.EmbeddedAttachments) > 0 {
			err := multipartWriter.Close()
			if err != nil {
				return writer.Bytes(), err
			}

			lastBytes = writer.Bytes()
			lastBoundary = multipartWriter.Boundary()
			writer = &bytes.Buffer{}
			multipartWriter = multipart.NewWriter(writer)

			mimeHeader := textproto.MIMEHeader{}

			mimeHeader.Add("Content-Type", fmt.Sprintf(`%s;%s boundary="%s"`, MultipartAlternative, Seperator, lastBoundary))

			w, err := multipartWriter.CreatePart(mimeHeader)
			if err != nil {
				return writer.Bytes(), err
			}

			_, err = w.Write(lastBytes)
			if err != nil {
				return writer.Bytes(), err
			}
		}
	}

	if len(self.EmbeddedAttachments) > 0 {
		self.Header.ContentType = MultipartRelated

		for _, v := range self.EmbeddedAttachments {
			mimeHeader := textproto.MIMEHeader{}

			mimeHeader.Add("Content-Type", fmt.Sprintf(`%s;%s name="%s"`, v.ContentType, Seperator, v.Name))
			mimeHeader.Add("Content-Transfer-Encoding", v.ContentTransferEncoding)
			mimeHeader.Add("Content-ID", fmt.Sprintf("<%s>", v.ContentId))

			w, err := multipartWriter.CreatePart(mimeHeader)
			if err != nil {
				continue
			}

			_, err = w.Write(v.Data)
			if err != nil {
				continue
			}
		}

		if len(self.Attachments) > 0 {
			err := multipartWriter.Close()
			if err != nil {
				return writer.Bytes(), err
			}

			lastBytes = writer.Bytes()
			lastBoundary = multipartWriter.Boundary()
			writer = &bytes.Buffer{}
			multipartWriter = multipart.NewWriter(writer)

			mimeHeader := textproto.MIMEHeader{}

			mimeHeader.Add("Content-Type", fmt.Sprintf(`%s;%s boundary="%s"`, MultipartRelated, Seperator, lastBoundary))

			w, err := multipartWriter.CreatePart(mimeHeader)
			if err != nil {
				return writer.Bytes(), err
			}

			_, err = w.Write(lastBytes)
			if err != nil {
				return writer.Bytes(), err
			}
		}
	}

	if len(self.Attachments) > 0 {
		self.Header.ContentType = MultipartMixed

		for _, v := range self.Attachments {
			mimeHeader := textproto.MIMEHeader{}

			mimeHeader.Add("Content-Type", fmt.Sprintf(`%s;%s name="%s"`, v.ContentType, Seperator, mime.BEncoding.Encode(DefaultEncoding, v.Name)))
			mimeHeader.Add("Content-Transfer-Encoding", v.ContentTransferEncoding)
			mimeHeader.Add("Content-Disposition", fmt.Sprintf(`%s;%s filename="%s"`, DispositionAttachment, Seperator, mime.BEncoding.Encode(DefaultEncoding, v.FileName)))

			w, err := multipartWriter.CreatePart(mimeHeader)
			if err != nil {
				continue
			}

			_, err = w.Write(v.Data)
			if err != nil {
				continue
			}
		}
	}

	err := multipartWriter.Close()
	if err != nil {
		return writer.Bytes(), err
	}

	lastBytes = writer.Bytes()
	lastBoundary = multipartWriter.Boundary()
	writer = &bytes.Buffer{}
	multipartWriter = multipart.NewWriter(writer)

	err = self.Header.Write(writer, lastBoundary)
	if err != nil {
		return writer.Bytes(), err
	}

	_, err = writer.Write(lastBytes)
	if err != nil {
		return writer.Bytes(), err
	}

	return writer.Bytes(), nil
}
