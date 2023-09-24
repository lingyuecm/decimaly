package model

import "testing"

func TestCreateNumber(t *testing.T) {
	n, e := CreateNumber("+13413.13241")
	t.Log(e)
	t.Log(n)
	_, e = CreateNumber("123e4")
	t.Log(e)
}

func TestNumber_Add(t *testing.T) {
	n1, _ := CreateNumber("+13413.123")
	n2, _ := CreateNumber("-13413")

	n := n1.Add(n2)
	t.Log(n.scale)
	t.Log(n.digits)
}
