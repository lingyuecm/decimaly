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
	sa := append(make([]Segment, 1, 1), shrinkUnsigned(unsignedMultiplication(sa1, sa2))...)
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
	q = append(make([]Segment, 1, 1), q...)
	if sign > 0 {
		q = generateNegative(q)
	}
	r = append(make([]Segment, 1, 1), r...)
	if b1.complement[0] > 0 {
		r = generateNegative(r)
	}
	quotient := new(BinaryInteger)
	quotient.complement = q

	remainder := new(BinaryInteger)
	remainder.complement = r

	return quotient, remainder, nil
}

func complementAddition(sa1 []Segment, sa2 []Segment) []Segment {
	l1 := len(sa1)
	l2 := len(sa2)
	l := bigger(l1, l2)

	var s1 Segment
	var s2 Segment
	var index1 int
	var index2 int

	expandedSign1 := expandSign(sa1[0])
	expandedSign2 := expandSign(sa2[0])
	result := make([]Segment, l, l)
	var sum Segment
	var carry DoubleSegment = 0
	for m := 1; m <= l; m++ {
		index1 = l1 - m
		if index1 > 0 {
			s1 = sa1[index1]
		} else {
			s1 = expandedSign1
		}

		index2 = l2 - m
		if index2 > 0 {
			s2 = sa2[index2]
		} else {
			s2 = expandedSign2
		}
		sum, carry = segmentAddition(s1, s2, carry)
		result[l-m] = sum
	}
	return result
}

func adjustComplement(sa []Segment) []Segment {
	result := make([]Segment, 1, 1)
	result[0] = sa[0] >> (segmentLength - 1)
	expandedSign := expandSign(result[0])
	length := len(sa)
	for m := 0; m < length; m++ {
		if sa[m] != expandedSign {
			return append(result, sa[m:]...)
		}
	}
	return append(result, expandedSign)
}

func shrinkUnsigned(sa []Segment) []Segment {
	if sa[0] > 0 {
		return sa
	}
	length := len(sa)
	for m := 1; m < length; m++ {
		if sa[m] > 0 {
			return sa[m:]
		}
	}
	return make([]Segment, 1, 1) // 0
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

func unsignedAddition(sa1 []Segment, sa2 []Segment) []Segment {
	l1 := len(sa1)
	l2 := len(sa2)
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
			s1 = sa1[index1]
		} else {
			s1 = 0
		}

		index2 = l2 - m
		if index2 >= 0 {
			s2 = sa2[index2]
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

func unsignedMultiplication(sa1 []Segment, sa2 []Segment) []Segment {
	l2 := len(sa2)
	result := make([]Segment, 0, 0)
	for m := 0; m < l2; m++ {
		result = unsignedAddition(shiftSegmentL(result, 1), generatePartialProduct(sa1, sa2[m]))
	}
	return result
}

func unsignedDivision(sa1 []Segment, sa2 []Segment) ([]Segment, []Segment) { // Quotient, Remainder
	c := unsignedComparison(sa1, sa2)
	if c < 0 {
		return make([]Segment, 1, 1), sa1
	} else if c == 0 {
		q := make([]Segment, 1, 1)
		q[0] = 1
		return q, make([]Segment, 1, 1)
	}

	l1 := len(sa1)
	l2 := len(sa2)
	l := l1 - l2 + 1
	result := make([]Segment, l, l)

	index1 := l2 - 1
	index := 0
	var q Segment
	r := sa1[0:index1]
	for {
		if index1 >= l1 {
			break
		}
		q, r = findPartialQuotient(append(r, sa1[index1]), sa2)
		if len(r) == 1 && r[0] == 0 {
			r = make([]Segment, 0, 0)
		}
		result[index] = q
		index++

		index1++
	}
	if len(r) == 0 {
		r = make([]Segment, 1, 1)
	}
	return shrinkUnsigned(result[0:index]), r
}

func unsignedComparison(sa1 []Segment, sa2 []Segment) int64 {
	l1 := len(sa1)
	l2 := len(sa2)
	if l1 > l2 {
		return 1
	} else if l1 < l2 {
		return -1
	}
	// l1 == l2
	for m := 0; m < l1; m++ {
		if sa1[m] != sa2[m] {
			return int64(sa1[m]) - int64(sa2[m])
		}
	}
	return 0
}

func findPartialQuotient(sa1 []Segment, sa2 []Segment) (Segment, []Segment) { // The Segment of the Quotient, Remainder
	c := unsignedComparison(sa1, sa2)
	if c < 0 {
		return 0, sa1
	} else if c == 0 {
		return 1, make([]Segment, 1, 1)
	}

	var q Segment = 0
	var qb Segment
	var r = sa1
	for {
		qb, r = findPartialQuotientBit(r, sa2)
		q = q + qb
		if unsignedComparison(r, sa2) < 0 {
			return q, r
		}
	}
}

func findPartialQuotientBit(sa1 []Segment, sa2 []Segment) (Segment, []Segment) { // The Bit of the Segment, Remainder
	bl1 := unsignedBitLength(sa1)
	bl2 := unsignedBitLength(sa2)
	q := shiftBitsL(sa2, bl1-bl2)
	c := unsignedComparison(sa1, q)
	if c >= 0 {
		return 1 << (bl1 - bl2), unsignedSubtraction(sa1, q)
	}
	return 1 << (bl1 - bl2 - 1), unsignedSubtraction(sa1, shiftBitR(q))
}

func unsignedSubtraction(sa1 []Segment, sa2 []Segment) []Segment {
	l1 := len(sa1)
	l2 := len(sa2)
	l := bigger(l1, l2)
	result := make([]Segment, l, l)
	var carry DoubleSegment = 0
	var s1 Segment
	var s2 Segment
	var index1 int
	var index2 int
	for m := 1; m <= l; m++ {
		index1 = l1 - m
		if index1 >= 0 {
			s1 = sa1[index1]
		} else {
			s1 = 0
		}
		index2 = l2 - m
		if index2 >= 0 {
			s2 = sa2[index2]
		} else {
			s2 = 0
		}
		result[l-m], carry = segmentSubtraction(s1, s2, carry)
	}
	return shrinkUnsigned(result)
}

func unsignedBitLength(sa []Segment) uint64 {
	s := sa[0]
	if 0 == s {
		return 0
	}

	l := segmentLength
	ds := DoubleSegment(s)
	m := 0
	for {
		ds = ds << 1
		if ds&carryThreshold > 0 {
			return uint64(l-m) + uint64(len(sa)-1)*uint64(segmentLength)
		}
		m = m + 1
	}
}

func shiftBitsL(sa []Segment, bitCount uint64) []Segment {
	sl := uint64(segmentLength)
	if bitCount == 0 {
		return sa
	} else if bitCount < sl {
		var factor Segment = 1 << bitCount
		return generatePartialProduct(sa, factor)
	} else {
		return shiftBitsL(shiftSegmentL(sa, int(bitCount/sl)), bitCount%sl)
	}
}

func shiftBitR(sa []Segment) []Segment {
	l := len(sa)
	result := make([]Segment, l, l)
	var r Segment = 0
	for m := 0; m < l; m++ {
		result[m] = (r << (segmentLength - 1)) + sa[m]/2
		r = sa[m] % 2
	}
	return shrinkUnsigned(result)
}

func generatePartialProduct(sa1 []Segment, sa2 Segment) []Segment {
	length := len(sa1)
	result := make([]Segment, length+1, length+1)

	var product Segment
	var carry DoubleSegment = 0

	for m := length; m >= 1; m-- {
		product, carry = segmentMultiplication(sa1[m-1], sa2, carry)
		result[m] = product
	}
	if carry > 0 {
		result[0] = Segment(carry)
		return result
	}
	return result[1:]
}

func generateComplement(sa []Segment, sign Segment) []Segment {
	if signPositive == sign {
		return sa
	}
	return generateNegative(sa)
}

func generateNegative(sa []Segment) []Segment {
	length := len(sa)
	sum := make([]Segment, length, length)
	sum[0] = 2
	result := make([]Segment, length, length)
	var difference Segment
	var carry DoubleSegment = 0
	for m := length - 1; m >= 0; m-- {
		difference, carry = segmentSubtraction(sum[m], sa[m], carry)
		result[m] = difference
	}
	result[0] = result[0] & 1 // -0
	return result
}

func expandSign(sign Segment) Segment {
	return Segment((carryThreshold - DoubleSegment(sign)) % carryThreshold)
}

func shiftSegmentL(sa []Segment, count int) []Segment {
	return append(sa, make([]Segment, count, count)...)
}
