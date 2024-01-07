package tools

import (
	"testing"
	"strings"
)

func TestEncryptDecrypt(t *testing.T) {
	// Given
	str := "test string"
	pass := "1234"
	encCmd := NewEncryptArmoredStdOutCmd(pass)
	decCmd := NewDecrypt(pass)

	// When
	encCmd.Stdin = strings.NewReader(str)
	enc, _ := encCmd.Output()

	decCmd.Stdin = strings.NewReader(string(enc))
	dec, _ := decCmd.Output()

	// Then
	if string(dec) != str  {
		t.Errorf("Expected '%s', but got '%s'", str, dec)
	}
}
