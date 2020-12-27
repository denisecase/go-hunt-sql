package formaterror

import (
	"errors"
	"strings"
)

// FormatError explains user input issues
func FormatError(err string) error {

	if strings.Contains(err, "email") {
		return errors.New("Email Already Taken")
	}

	if strings.Contains(err, "name") {
		return errors.New("Name Already Taken")
	}
	if strings.Contains(err, "hashedPassword") {
		return errors.New("Incorrect Password")
	}
	return errors.New("Incorrect Details")
}
