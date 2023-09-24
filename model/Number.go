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
