package tools

import (
	"os/exec"
	"strings"
)

const (
	tar       = "tar"
	gpg       = "gpg"
	separator = " "
	required  = tar + separator + gpg
)

func checkIfToolExists(prog string) error {
	_, err := exec.LookPath(prog)
	return err
}

func CheckRequiredTools() error {
	for _, element := range strings.Split(required, separator) {
		if err := checkIfToolExists(element); err != nil {
			return err
		}
	}

	return nil
}
