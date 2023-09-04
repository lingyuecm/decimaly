package model

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"unsafe"
)

type Segment = uint8
type DoubleSegment = uint16

const segmentLength = 8

const segmentMask DoubleSegment = 1<<segmentLength - 1
const carryThreshold DoubleSegment = 1 << segmentLength
const ten Segment = 10

const signPositive Segment = 0
const signNegative Segment = 1
const maxIndex10 = "10000000000000000000"

var expandFactor = math.Log(float64(carryThreshold)) / math.Log(float64(ten))

type BinaryInteger struct {
	complement []Segment
}

func CreateBinaryInteger(value string) (*BinaryInteger, error) {
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

	result, err := createUnsigned(value[startIndex:])
	if nil != err {
		return nil, err
	}
	result = append(make([]Segment, 1), result...)

	i := new(BinaryInteger)
	i.complement = generateComplement(result, sign)
	return i, nil
}

func (b1 *BinaryInteger) Negative() *BinaryInteger {
	i := new(BinaryInteger)
	i.complement = generateNegative(b1.complement)
	return i
}

func (b1 *BinaryInteger) Add(b2 *BinaryInteger) *BinaryInteger {
	i := new(BinaryInteger)
	i.complement = adjustComplement(complementAddition(b1.complement, b2.complement))
	return i
}

func (b1 *BinaryInteger) Subtract(b2 *BinaryInteger) *BinaryInteger {
	return b1.Add(b2.Negative())
}

func (b1 *BinaryInteger) Multiply(b2 *BinaryInteger) *BinaryInteger {
	sa1 := b1.complement
	sa2 := b2.complement
	sign := (sa1[0] + sa2[0]) & 1

	if sa1[0] > 0 {
		sa1 = generateNegative(sa1)
	}
	if sa2[0] > 0 {
		sa2 = generateNegative(sa2)
	}
	sa := append(make([]Segment, 1), shrinkUnsigned(unsignedMultiplication(sa1, sa2))...)
	if sign > 0 {
		sa = generateNegative(sa)
	}

	i := new(BinaryInteger)
	i.complement = sa
	return i
}

func (b1 *BinaryInteger) DividedBy(b2 *BinaryInteger) (*BinaryInteger, *BinaryInteger, error) {
	if len(b2.complement) == 2 && b2.complement[0] == 0 && b2.complement[1] == 0 {
		return nil, nil, errors.New(fmt.Sprintf("Cannot Be Divided by Zero"))
	}

	sa1 := b1.complement
	sa2 := b2.complement
	sign := (sa1[0] + sa2[0]) & 1

	if sa1[0] > 0 {
		sa1 = generateNegative(sa1)
	}
	if sa2[0] > 0 {
		sa2 = generateNegative(sa2)
	}
	q, r := unsignedDivision(shrinkUnsigned(sa1), shrinkUnsigned(sa2))
	q = append(make([]Segment, 1), q...)
	if sign > 0 {
		q = generateNegative(q)
	}
	r = append(make([]Segment, 1), r...)
	if b1.complement[0] > 0 {
		r = generateNegative(r)
	}
	quotient := new(BinaryInteger)
	quotient.complement = q

	remainder := new(BinaryInteger)
	remainder.complement = r

	return quotient, remainder, nil
}

func (b1 *BinaryInteger) GcdWith(b2 *BinaryInteger) *BinaryInteger {
	sa1 := b1.complement
	if sa1[0] > 0 {
		sa1 = generateNegative(sa1)
	}
	sa1 = shrinkUnsigned(sa1)

	sa2 := b2.complement
	if sa2[0] > 0 {
		sa2 = generateNegative(sa2)
	}
	sa2 = shrinkUnsigned(sa2)

	_, r := unsignedDivision(sa1, sa2)
	for {
		if r[0] == 0 {
			i := new(BinaryInteger)
			i.complement = append(make([]Segment, 1), sa2...)
			return i
		}
		sa1 = sa2
		sa2 = r
		_, r = unsignedDivision(sa1, sa2)
	}
}

func (b1 *BinaryInteger) DecimalValue() string {
	sa := b1.complement
	sign := ""
	if sa[0] > 0 {
		sign = "-"
		sa = generateNegative(sa)
	}
	sa = shrinkUnsigned(sa)
	divider, _ := createUnsigned(maxIndex10)
	length := int(float64(len(sa))*expandFactor) + 2
	result := make([]byte, length)
	index := length
	var group string
	var gl int
	var q []Segment
	var r []Segment
	for {
		q, r = unsignedDivision(sa, divider)
		group = strconv.FormatUint(getUint64(r), 10)
		gl = len(group)
		for m := gl - 1; m >= 0; m-- {
			index--
			result[index] = group[m]
		}
		if q[0] == 0 {
			break
		}
		for m := gl; m < 19; m++ {
			index--
			result[index] = '0'
		}
		sa = q
	}
	if len(sign) > 0 {
		index--
		result[index] = sign[0]
	}
	result = result[index:]
	return *(*string)(unsafe.Pointer(&result))
}
