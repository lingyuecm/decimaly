package model

import (
	"errors"
	"fmt"
)

func createUnsigned(value string) ([]Segment, error) {
	return createUnsignedBasedOn(value, make([]Segment, 1))
}

func createUnsignedBasedOn(value string, baseValue []Segment) ([]Segment, error) {
	length := len(value)
	result := baseValue
	digit := make([]Segment, 1)
	var nonZeroIndex int
	for m := 0; m < length; m++ {
		if value[m] != '0' {
			nonZeroIndex = m
			break
		}
	}
	for m := nonZeroIndex; m < length; m++ {
		if value[m] < '0' || value[m] > '9' {
			return nil, errors.New(fmt.Sprintf("Invalid Digit at %d: %c", m, value[m]))
		}
		digit[0] = Segment(value[m] - '0')
		result = unsignedAddition(generatePartialProduct(result, ten), digit)
	}
	return result, nil
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
	result := make([]Segment, l)
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
	result := make([]Segment, 1)
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
	return make([]Segment, 1) // 0
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

	result := make([]Segment, l+1)
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
	result := make([]Segment, 0)
	for m := 0; m < l2; m++ {
		result = unsignedAddition(shiftSegmentL(result, 1), generatePartialProduct(sa1, sa2[m]))
	}
	return result
}

func unsignedDivision(sa1 []Segment, sa2 []Segment) ([]Segment, []Segment) { // Quotient, Remainder
	if len(sa1) < len(sa2) {
		return make([]Segment, 1), sa1
	}

	l1 := len(sa1)
	l2 := len(sa2)
	l := l1 - l2 + 1
	result := make([]Segment, l)

	index1 := l2 - 1
	index := 0
	var q Segment
	r := sa1[0:index1]
	for {
		if index1 >= l1 {
			break
		}
		q, r = findPartialQuotient(append(r, sa1[index1]), sa2)
		if r[0] == 0 {
			r = make([]Segment, 0)
		}
		result[index] = q
		index++

		index1++
	}
	if len(r) == 0 {
		r = make([]Segment, 1)
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
		return 1, make([]Segment, 1)
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
	result := make([]Segment, l)
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
	result := make([]Segment, l)
	var r Segment = 0
	for m := 0; m < l; m++ {
		result[m] = (r << (segmentLength - 1)) + sa[m]/2
		r = sa[m] % 2
	}
	return shrinkUnsigned(result)
}

func generatePartialProduct(sa1 []Segment, sa2 Segment) []Segment {
	length := len(sa1)
	result := make([]Segment, length+1)

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
	sum := make([]Segment, length)
	sum[0] = 2
	result := make([]Segment, length)
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
	return append(sa, make([]Segment, count)...)
}

func getUint64(sa []Segment) uint64 {
	carryThreshold64 := uint64(carryThreshold)
	var sum uint64 = 0
	for _, s := range sa {
		sum = sum*carryThreshold64 + uint64(s)
	}
	return sum
}
