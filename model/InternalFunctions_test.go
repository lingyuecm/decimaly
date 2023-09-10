package model

import "testing"

func TestCreateUnsigned(t *testing.T) {
	t.Log(createUnsigned("00000000000000000000001234"))
}
