package model

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type Number struct {
	digits *BinaryInteger
	scale  int
}

func CreateNumber(value string) (*Number, error) {
	if r, _ := regexp.Compile("^[+\\-]?\\d+(\\.\\d+)?$"); !r.MatchString(value) {
		return nil, errors.New(fmt.Sprintf("Invalid Number Format"))
	}
	length := len(value)
	dotIndex := strings.Index(value, ".")
	if dotIndex < 0 {
		n := new(Number)

		n.digits, _ = CreateBinaryInteger(value)
		n.scale = 0

		return n, nil
	}

	var startIndex = 0
	var sign = signPositive
	if '+' == value[0] {
		startIndex = 1
	} else if '-' == value[0] {
		startIndex = 1
		sign = signNegative
	}

	segments, _ := createUnsigned(value[startIndex:dotIndex])
	segments, _ = createUnsignedBasedOn(value[dotIndex+1:], segments)

	i := new(BinaryInteger)
	i.sign = sign
	i.segments = segments

	n := new(Number)
	n.digits = i
	n.scale = length - dotIndex - 1

	return n, nil
}

func (n1 *Number) Add(n2 *Number) *Number {
	tenBin, _ := CreateBinaryInteger("10")

	b1 := n1.digits
	b2 := n2.digits

	if n1.scale < n2.scale {
		for m := n1.scale; m < n2.scale; m++ {
			b1 = b1.Multiply(tenBin)
		}
	} else if n1.scale > n2.scale {
		for m := n2.scale; m < n1.scale; m++ {
			b2 = b2.Multiply(tenBin)
		}
	}

	n := new(Number)

	n.digits = b1.Add(b2)
	n.scale = bigger(n1.scale, n2.scale)

	return n
}
