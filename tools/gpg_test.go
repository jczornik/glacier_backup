package tools

import (
	"strings"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	// Given
	str := "test string"
	pass := "1234"
	encCmd := NewEncryptArmoredStdOutCmd(pass)
	decCmd := NewDecrypt(pass)

	// When
	encCmd.cmd.Stdin = strings.NewReader(str)
	enc, _ := encCmd.cmd.Output()

	decCmd.cmd.Stdin = strings.NewReader(string(enc))
	dec, _ := decCmd.cmd.Output()

	// Then
	if string(dec) != str {
		t.Errorf("Expected '%s', but got '%s'", str, dec)
	}
}
