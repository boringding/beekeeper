//Helper functions for mime handling.

package mailmsg

import (
	"errors"
	"fmt"
	"io"
	"net/textproto"
)

const (
	Seperator = "\r\n"
)

func writeHeader(header textproto.MIMEHeader, writer io.Writer) error {
	for k, v := range header {
		if len(v) != 1 {
			return errors.New("value count incorrect")
		}

		_, err := fmt.Fprintf(writer, "%s: %s%s", k, textproto.TrimString(v[0]), Seperator)
		if err != nil {
			return err
		}
	}

	_, err := fmt.Fprintf(writer, "%s", Seperator)
	if err != nil {
		return err
	}

	return nil
}
