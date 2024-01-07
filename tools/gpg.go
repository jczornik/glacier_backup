package tools

import "os/exec"

func NewEncryptArmoredStdOutCmd(pass string) *exec.Cmd {
	return exec.Command(gpg, "-c", "--cipher-algo", "AES256", "--batch", "--passphrase", pass)
}

func NewDecrypt(pass string) *exec.Cmd {
	return exec.Command(gpg, "-d", "--batch", "--passphrase", pass)
}
