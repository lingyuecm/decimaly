package model

import (
	"math/rand"
	"testing"
	"unsafe"
)

func TestCreateBinaryInteger(t *testing.T) {
	l := 30000
	value := make([]byte, l)
	value[0] = '9'
	for m := 1; m < l; m++ {
		value[m] = '0' + byte(rand.Intn(10))
	}
	valueText := *(*string)(unsafe.Pointer(&value))
	b, _ := CreateBinaryInteger(valueText)
	t.Log(nil == b)
	t.Log(b.DecimalValue() == valueText)
}

func TestBinaryInteger_Negative(t *testing.T) {
	b, _ := CreateBinaryInteger("123")
	t.Log(b.Negative().DecimalValue())

	b2, _ := CreateBinaryInteger("0")
	t.Log(b2.Negative().DecimalValue())
}

func TestBinaryInteger_Add(t *testing.T) {
	b1, _ := CreateBinaryInteger("-30498579023579024377")
	b2, _ := CreateBinaryInteger("30498579023579024376")
	b3 := b1.Add(b2)
	t.Log(b3.DecimalValue())
}

func TestBinaryInteger_Subtract(t *testing.T) {
	b1, _ := CreateBinaryInteger("123")
	b2, _ := CreateBinaryInteger("4567")
	b3 := b1.Subtract(b2)
	t.Log(b3.DecimalValue())
}

func TestBinaryInteger_Multiply(t *testing.T) {
	b1, _ := CreateBinaryInteger("123")
	b2, _ := CreateBinaryInteger("4567")
	b3 := b1.Multiply(b2)
	t.Log(b3.DecimalValue())
}

func TestBinaryInteger_DividedBy(t *testing.T) {
	b1, _ := CreateBinaryInteger("123")
	b2, _ := CreateBinaryInteger("4567")
	b3, b4, _ := b1.DividedBy(b2)
	t.Log(b3.DecimalValue())
	t.Log(b4.DecimalValue())
}

func TestBinaryInteger_GcdWith(t *testing.T) {
	b1, _ := CreateBinaryInteger("-123")
	b2, _ := CreateBinaryInteger("82")
	b3 := b1.GcdWith(b2)
	t.Log(b3.DecimalValue())
}
