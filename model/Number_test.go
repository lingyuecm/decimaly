package model

import "testing"

func TestCreateNumber(t *testing.T) {
	n, e := CreateNumber("+13413.13241")
	t.Log(e)
	t.Log(n)
	_, e = CreateNumber("123e4")
	t.Log(e)
}

func TestNumber_Negative(t *testing.T) {
	n, e := CreateNumber("+13413.13241")
	t.Log(e)
	t.Log(n.DecimalValue())
	t.Log(n.Negative().DecimalValue())

	n1, _ := CreateNumber("-0")
	t.Log(n1.Negative().DecimalValue())
}

func TestNumber_Add(t *testing.T) {
	n1, _ := CreateNumber("+13413.123")
	n2, _ := CreateNumber("-13413")

	n := n1.Add(n2)
	t.Log(n.DecimalValue())

	n3, _ := CreateNumber("-13413.123")
	t.Log(n1.Add(n3).DecimalValue())
	t.Log(n3.Add(n1).DecimalValue())
}

func TestNumber_Subtraction(t *testing.T) {
	n1, _ := CreateNumber("2139847.9803745289075")
	n2, _ := CreateNumber("+2139847.9803745289075")
	t.Log(n1.Subtraction(n2).DecimalValue())

	n3, _ := CreateNumber("7890367896789523.12378904059103")
	t.Log(n1.Subtraction(n3).DecimalValue())
}

func TestNumber_DecimalValue(t *testing.T) {
	n1, _ := CreateNumber("+13413.123")
	t.Log(n1.DecimalValue())

	n2, _ := CreateNumber("0.12345")
	t.Log(n2.DecimalValue())

	n3, _ := CreateNumber("-0.012345")
	t.Log(n3.DecimalValue())

	n4, _ := CreateNumber("+0.0012345")
	t.Log(n4.DecimalValue())

	n5, _ := CreateNumber("-1237645123784.0000000000123749102347910")
	t.Log(n5.DecimalValue())

	n6, _ := CreateNumber("1237645123784.123749102347910")
	t.Log(n6.DecimalValue())

	n7, _ := CreateNumber("-384567234895676")
	t.Log(n7.DecimalValue())

	n8, _ := CreateNumber("384567234895676")
	t.Log(n8.DecimalValue())
}
