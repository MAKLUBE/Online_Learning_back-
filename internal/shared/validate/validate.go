package validate

import (
	"net/mail"
	"strings"
)

func Required(value string) bool {
	return strings.TrimSpace(value) != ""
}

func Email(value string) bool {
	address, err := mail.ParseAddress(strings.TrimSpace(value))
	return err == nil && address.Address == strings.TrimSpace(value)
}

func MinLen(value string, min int) bool {
	return len(strings.TrimSpace(value)) >= min
}
