package model

import (
	"errors"
	"fmt"
)

type Segment = uint8
type DoubleSegment = uint16

const segmentLength = 8

const segmentMask DoubleSegment = 1<<segmentLength - 1
const carryThreshold DoubleSegment = 1 << segmentLength
const ten Segment = 10

const signPositive Segment = 0
const signNegative Segment = 1

type BinaryInteger struct {
	complement []Segment
}

func CreateBigInteger(value string) (*BinaryInteger, error) {
	length := len(value)
	if 0 == length {
		return nil, errors.New(fmt.Sprintf("Empty String"))
	}
	var startIndex = 0
	var sign = signPositive
	if '+' == value[0] {
		startIndex = 1
	} else if '-' == value[0] {
		startIndex = 1
		sign = signNegative
	}

	result := make([]Segment, 1, 1)
	digit := make([]Segment, 1, 1)
	var found = false
	for m := startIndex; m < length; m++ {
		if value[m] < '0' || value[m] > '9' {
			return nil, errors.New(fmt.Sprintf("Invalid Digit at %d: %c", m, value[m]))
		}
		if '0' == value[m] && !found {
			continue
		} else {
			found = true
		}
		digit[0] = value[m] - '0'
		result = unsignedAddition(generatePartialProduct(result, ten), digit)
	}
	result = append(make([]Segment, 1, 1), result...)

	i := new(BinaryInteger)
	i.complement = generateComplement(result, sign)
	return i, nil
}

func (b *BinaryInteger) Negative() *BinaryInteger {
	i := new(BinaryInteger)
	i.complement = generateNegative(b.complement)
	return i
}

func segmentAddition(s1 Segment, s2 Segment, carry DoubleSegment) (Segment, DoubleSegment) { // Sum, Carry
	sum := DoubleSegment(s1) + DoubleSegment(s2) + carry
	return Segment(sum & segmentMask), (sum >> segmentLength) & segmentMask
}

func segmentSubtraction(s1 Segment, s2 Segment, carry DoubleSegment) (Segment, DoubleSegment) { // Difference, Carry
	difference := DoubleSegment(s1) + carryThreshold - DoubleSegment(s2) - carry
	return Segment(difference % carryThreshold), 1 - difference/carryThreshold
}

func segmentMultiplication(s1 Segment, s2 Segment, carry DoubleSegment) (Segment, DoubleSegment) { // Product, Carry
	product := DoubleSegment(s1)*DoubleSegment(s2) + carry
	return Segment(product & segmentMask), (product >> segmentLength) & segmentMask
}

func bigger(a1 int, a2 int) int {
	if a1 >= a2 {
		return a1
	}
	return a2
}

func unsignedAddition(b1 []Segment, b2 []Segment) []Segment {
	l1 := len(b1)
	l2 := len(b2)
	l := bigger(l1, l2)

	var s1 Segment
	var s2 Segment
	var index1 int
	var index2 int

	result := make([]Segment, l+1, l+1)
	var sum Segment
	var carry DoubleSegment = 0
	for m := 1; m <= l; m++ {
		index1 = l1 - m
		if index1 >= 0 {
			s1 = b1[index1]
		} else {
			s1 = 0
		}

		index2 = l2 - m
		if index2 >= 0 {
			s2 = b2[index2]
		} else {
			s2 = 0
		}
		sum, carry = segmentAddition(s1, s2, carry)
		result[l+1-m] = sum
	}
	if carry > 0 {
		result[0] = Segment(carry)
		return result
	}
	return result[1:]
}

func unsignedMultiplication(s1 []Segment, s2 []Segment) []Segment {
	l2 := len(s2)
	result := make([]Segment, 0, 0)
	for m := 0; m < l2; m++ {
		result = unsignedAddition(shiftSegmentL(result, 1), generatePartialProduct(s1, s2[m]))
	}
	return result
}

func generatePartialProduct(s1 []Segment, s2 Segment) []Segment {
	length := len(s1)
	result := make([]Segment, length+1, length+1)

	var product Segment
	var carry DoubleSegment = 0

	for m := length; m >= 1; m-- {
		product, carry = segmentMultiplication(s1[m-1], s2, carry)
		result[m] = product
	}
	if carry > 0 {
		result[0] = Segment(carry)
		return result
	}
	return result[1:]
}

func generateComplement(s []Segment, sign Segment) []Segment {
	if signPositive == sign {
		return s
	}
	return generateNegative(s)
}

func generateNegative(s []Segment) []Segment {
	length := len(s)
	sum := make([]Segment, length, length)
	sum[0] = 2
	result := make([]Segment, length, length)
	var difference Segment
	var carry DoubleSegment = 0
	for m := length - 1; m >= 0; m-- {
		difference, carry = segmentSubtraction(sum[m], s[m], carry)
		result[m] = difference
	}
	return result
}

func shiftSegmentL(s []Segment, count int) []Segment {
	return append(s, make([]Segment, count, count)...)
}
