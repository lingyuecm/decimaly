package model

import "testing"

func TestCreateNumber(t *testing.T) {
	n, e := CreateNumber("+13413.13241")
	t.Log(e)
	t.Log(n)
	_, e = CreateNumber("123e4")
	t.Log(e)
}
