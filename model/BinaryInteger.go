package model

import (
	"errors"
	"fmt"
	"strconv"
	"unsafe"
)

// BinaryInteger is the binary representation of integers
// Considering there isn't a kind of data type consisting of 1 bit, to make full use of the memory,
// certain count (say, s) of bits will be grouped as a Segment, which is an alias of certain unsigned integer types.
// So such an integer is actually one of radix 2^s
type BinaryInteger struct {
	sign     Segment   // The Sign
	segments []Segment // The Bits Grouped as Segments
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

	i := new(BinaryInteger)
	i.sign = sign
	i.segments = result

	return i, nil
}

func (b1 *BinaryInteger) IsZero() bool {
	return len(b1.segments) == 1 && b1.segments[0] == 0
}

func (b1 *BinaryInteger) Negative() *BinaryInteger {
	i := new(BinaryInteger)

	i.sign = signNegative - b1.sign
	i.segments = b1.segments

	return i
}

func (b1 *BinaryInteger) Add(b2 *BinaryInteger) *BinaryInteger {
	sa1 := generateComplement(append(make([]Segment, 1), b1.segments...), b1.sign)
	sa2 := generateComplement(append(make([]Segment, 1), b2.segments...), b2.sign)

	complement := adjustComplement(complementAddition(sa1, sa2))

	i := new(BinaryInteger)
	i.sign = complement[0]
	if i.sign > 0 {
		i.segments = shrinkUnsigned(generateNegative(complement))
	} else {
		i.segments = shrinkUnsigned(complement)
	}
	return i
}

func (b1 *BinaryInteger) Subtract(b2 *BinaryInteger) *BinaryInteger {
	return b1.Add(b2.Negative())
}

func (b1 *BinaryInteger) Multiply(b2 *BinaryInteger) *BinaryInteger {
	i := new(BinaryInteger)

	i.sign = (b1.sign + b2.sign) & signNegative
	i.segments = shrinkUnsigned(unsignedMultiplication(b1.segments, b2.segments))

	return i
}

func (b1 *BinaryInteger) DividedBy(b2 *BinaryInteger) (*BinaryInteger, *BinaryInteger, error) {
	if b2.IsZero() {
		return nil, nil, errors.New(fmt.Sprintf("Cannot Be Divided by Zero"))
	}

	sign := (b1.sign + b2.sign) & signNegative
	q, r := unsignedDivision(b1.segments, b2.segments)

	quotient := new(BinaryInteger)
	quotient.sign = sign
	quotient.segments = q

	remainder := new(BinaryInteger)
	remainder.sign = b1.sign
	remainder.segments = r

	return quotient, remainder, nil
}

func (b1 *BinaryInteger) GcdWith(b2 *BinaryInteger) *BinaryInteger {
	sa1 := b1.segments
	sa2 := b2.segments

	_, r := unsignedDivision(sa1, sa2)
	for {
		if r[0] == 0 {
			i := new(BinaryInteger)
			i.segments = sa2
			return i
		}
		sa1 = sa2
		sa2 = r
		_, r = unsignedDivision(sa1, sa2)
	}
}

func (b1 *BinaryInteger) DecimalValue() string {
	if b1.IsZero() {
		return "0"
	}

	sa := b1.segments
	sign := ""
	if b1.sign > 0 {
		sign = "-"
	}
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
