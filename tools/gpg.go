package tools

import "os/exec"

func NewEncryptArmoredStdOutCmd(pass string) Cmd {
	cmd := exec.Command(gpg, "-c", "--cipher-algo", "AES256", "--batch", "--passphrase", pass)
	return NewCmd(cmd, nil)
}

func NewDecrypt(pass string) Cmd {
	cmd := exec.Command(gpg, "-d", "--batch", "--passphrase", pass)
	return NewCmd(cmd, nil)
}
