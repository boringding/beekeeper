package mailmsg

import (
	"errors"
	"net/mail"
	"strings"
)

const (
	AddrSeperator = ",\r\n "
)

func checkAddrs(addrs ...mail.Address) error {
	for _, v := range addrs {
		if v.Address == "" {
			return errors.New("find empty address")
		}
	}

	return nil
}

func stringAddrs(addrs ...mail.Address) string {
	addrStrs := make([]string, 0, 10)

	for _, v := range addrs {
		addrStrs = append(addrStrs, v.String())
	}

	return strings.Join(addrStrs, AddrSeperator)
}
