package model

import (
	"errors"
	"fmt"
	"unsafe"
)

const ten string = "1010"
const (
	signPositive byte = iota
	signNegative
)

var digitBinaries = [...]string{"0", "1", "10", "11", "100", "101", "110", "111", "1000", "1001"}

type BinaryInteger struct {
	complement string
}

func CreateBinaryInteger(value string) (*BinaryInteger, error) {
	length := len(value)
	if 0 == length {
		return nil, errors.New(fmt.Sprintf("Empty String"))
	}
	var sign byte
	var startIndex int
	if '+' == value[0] {
		sign = signPositive
		startIndex = 1
	} else if '-' == value[0] {
		sign = signNegative
		startIndex = 1
	} else {
		sign = signPositive
		startIndex = 0
	}

	original := "0"
	var found = false
	for m := startIndex; m < length; m++ {
		if value[m] > '9' || value[m] < '0' {
			return nil, errors.New(fmt.Sprintf("Invalid Digit at %d: %c", m, value[m]))
		}
		if '0' == value[m] && !found {
			continue
		} else {
			found = true
		}
		original = unsignedBinaryAddition(unsignedBinaryMultiplication(original, ten), digitBinaries[value[m]-'0'])
	}
	i := new(BinaryInteger)
	i.complement = shrink(generateComplement("0"+original, sign))

	return i, nil
}

func (b1 *BinaryInteger) Negative() *BinaryInteger {
	i := new(BinaryInteger)
	i.complement = generateNegative(b1.complement)
	return i
}

func (b1 *BinaryInteger) Abs() string {
	if signNegative == getSign(b1.complement) {
		return b1.Negative().Abs()
	}
	m := make(map[string]byte)
	for index, item := range digitBinaries {
		m[item] = '0' + byte(index)
	}

	complement := b1.complement[1:]
	// 2^3.3 Is about 9.8 so That the Number of Digits of Binaries Is about 3.3 Times of That of Decimals on Average
	decLength := int(float64(len(complement))/3.3) + 1 // Add 1 in Case of Tolerance
	decBytes := make([]byte, decLength, decLength)
	var q string
	var r string
	var index int
	for n := decLength - 1; n >= 0; n-- {
		index = n
		q, r = unsignedBinaryDivision(complement, ten)
		decBytes[n] = m[r]
		if "0" == q {
			break
		}
		complement = q
	}
	result := decBytes[index:]
	return *(*string)(unsafe.Pointer(&result))
}

func (b1 *BinaryInteger) Add(b2 *BinaryInteger) *BinaryInteger {
	result := new(BinaryInteger)
	result.complement = shrink(complementBinaryAddition(b1.complement, b2.complement))

	return result
}

func (b1 *BinaryInteger) Subtract(b2 *BinaryInteger) *BinaryInteger {
	return b1.Add(b2.Negative())
}

func (b1 *BinaryInteger) Multiply(b2 *BinaryInteger) *BinaryInteger {
	result := new(BinaryInteger)
	result.complement = complementBinaryMultiplication(b1.complement, b2.complement)
	return result
}

func generateComplement(original string, sign byte) string {
	if signPositive == sign {
		return original
	}
	return generateNegative(original)
}

func generateNegative(binaryCode string) string {
	length := len(binaryCode)
	bits := make([]byte, length, length)
	var diff byte
	var carry byte = 0
	for m := length - 1; m >= 0; m-- {
		diff, carry = bitwiseSubtraction(0, binaryCode[m]-'0', carry)
		bits[m] = '0' + diff
	}
	return *(*string)(unsafe.Pointer(&bits))
}

func bigger(a1 int, a2 int) int {
	if a1 >= a2 {
		return a1
	}
	return a2
}

func bitwiseAddition(bit1 byte, bit2 byte, carry byte) (byte, byte) { // Sum, Carry
	result := bit1 + bit2 + carry
	return result % 2, result / 2
}

func bitwiseSubtraction(bit1 byte, bit2 byte, carry byte) (byte, byte) {
	result := bit1 - bit2 - carry + 2
	return result % 2, 1 - result/2
}

func unsignedBinaryAddition(b1 string, b2 string) string {
	l1 := len(b1)
	l2 := len(b2)
	l := bigger(l1, l2)

	var index1 int
	var index2 int
	var bit1 byte
	var bit2 byte
	var sum byte
	var carry byte = 0

	result := make([]byte, l, l)

	for m := 1; m <= l; m++ {
		index1 = l1 - m
		if index1 >= 0 {
			bit1 = b1[index1] - '0'
		} else {
			bit1 = 0
		}
		index2 = l2 - m
		if index2 >= 0 {
			bit2 = b2[index2] - '0'
		} else {
			bit2 = 0
		}
		sum, carry = bitwiseAddition(bit1, bit2, carry)
		result[l-m] = '0' + sum
	}
	if carry > 0 {
		return "1" + *(*string)(unsafe.Pointer(&result))
	}
	return *(*string)(unsafe.Pointer(&result))
}

func complementBinaryAddition(b1 string, b2 string) string {
	l1 := len(b1)
	l2 := len(b2)
	l := bigger(l1, l2) + 1

	var index1 int
	var index2 int
	var bit1 byte
	var bit2 byte
	var sum byte
	var carry byte = 0

	result := make([]byte, l, l)

	for m := 1; m <= l; m++ {
		index1 = l1 - m
		if index1 >= 0 {
			bit1 = b1[index1] - '0'
		} else {
			bit1 = b1[0] - '0'
		}
		index2 = l2 - m
		if index2 >= 0 {
			bit2 = b2[index2] - '0'
		} else {
			bit2 = b2[0] - '0'
		}
		sum, carry = bitwiseAddition(bit1, bit2, carry)
		result[l-m] = '0' + sum
	}
	return *(*string)(unsafe.Pointer(&result))
}

func unsignedBinaryMultiplication(b1 string, b2 string) string {
	l2 := len(b2)

	result := ""

	for m := 0; m < l2; m++ {
		if '1' == b2[m] {
			result = unsignedBinaryAddition(result+"0", b1)
		} else {
			result = result + "0"
		}
	}
	return result
}

func complementBinaryMultiplication(b1 string, b2 string) string {
	l2 := len(b2)

	result := ""
	if signNegative == getSign(b2) {
		result = generateNegative(b1)
	}

	for m := 1; m < l2; m++ {
		if '1' == b2[m] {
			result = complementBinaryAddition(result+"0", b1)
		} else {
			result = result + "0"
		}
	}
	return shrink(result)
}

func unsignedBinaryDivision(b1 string, b2 string) (string, string) { // Quotient, Remainder
	c := unsignedBinaryComparison(b1, b2)
	if c < 0 {
		return "0", b1
	}

	l1 := len(b1)
	l2 := len(b2)
	if l1 == l2 {
		return "1", shrink0(complementBinaryAddition("0"+b1, generateNegative("0"+b2)))
	}
	q, r := findPartialQuotient(b1, b2)
	var l int
	qba := *(*[]byte)(unsafe.Pointer(&q)) // Quotient Byte Array
	ql := len(qba)
	for unsignedBinaryComparison(r, b2) >= 0 {
		l, r = findPartialQuotientLength(r, b2)
		qba[ql-l] = '1'
	}
	return *(*string)(unsafe.Pointer(&qba)), r
}

func findPartialQuotient(b1 string, b2 string) (string, string) { // Partial Quotient, Remainder
	l1 := len(b1)
	l2 := len(b2)
	l := l1 - l2
	q := shiftL("1", l)
	p := shiftL(b2, l)
	c := unsignedBinaryComparison(b1, p)
	if c < 0 {
		q = shiftR(q, 1)
		p = shiftR(p, 1)
	}
	return q, shrink0(complementBinaryAddition("0"+b1, generateNegative("0"+p)))
}

func findPartialQuotientLength(b1 string, b2 string) (int, string) { // The Length of the Partial Quotient, Remainder
	l1 := len(b1)
	l2 := len(b2)
	l := l1 - l2
	p := shiftL(b2, l)
	c := unsignedBinaryComparison(b1, p)
	if c < 0 {
		p = shiftR(p, 1)
		l = l - 1
	}
	return l + 1, shrink0(complementBinaryAddition("0"+b1, generateNegative("0"+p)))
}

func unsignedBinaryComparison(b1 string, b2 string) int {
	l1 := len(b1)
	l2 := len(b2)
	if l1 < l2 {
		return -1
	} else if l1 > l2 {
		return 1
	}
	for m := 0; m < l1; m++ {
		if b1[m] == b2[m] {
			continue
		}
		return int(b1[m]) - int(b2[m])
	}
	return 0
}

func shrink(complement string) string {
	length := len(complement)
	var found = false
	var index int
	for m := 0; m < length; m++ {
		if complement[m] != complement[0] {
			index = m
			found = true
			break
		}
	}
	if found {
		return complement[index-1:]
	}
	return complement[length-2:]
}

func shrink0(binary string) string {
	length := len(binary)
	for m := 0; m < length; m++ {
		if '1' == binary[m] {
			return binary[m:]
		}
	}
	return "0"
}

func getSign(complement string) byte {
	return complement[0] - '0'
}

func shiftL(b1 string, length int) string {
	zeros := make([]byte, length, length)
	for m := 0; m < length; m++ {
		zeros[m] = '0'
	}
	return b1 + *(*string)(unsafe.Pointer(&zeros))
}

func shiftR(b1 string, length int) string {
	l := len(b1)
	if length >= l {
		return "0"
	}
	return b1[0 : l-length]
}
