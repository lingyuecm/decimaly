package model

import "math"

/*
The most important constants of this framework are the type Segment, DoubleSegment and the constant segmentLength.
*/

// Segment is the "grouped bit" a BinaryInteger
/*
The reason why Segment is used instead of a certain unsigned integer type is that Go has several unsigned integer
types (uint8, uint16, uint32 and uint64), any one of which the user could choose to be used as a Segment and
the next one will be used as a DoubleSegment
*/
type Segment = uint8

// DoubleSegment is used for addition or multiplication where a carry may be generated
type DoubleSegment = uint16

// The length of a Segment
/*
Should be consistent with Segment. If Segment = uint8, the value of this constant should be 8; 16 for uint16 and so on
*/
const segmentLength = 8

/*
= = = = = = = = = = Derived constants = = = = = = = = = =
*/

// The mask to extract a Segment
const segmentMask DoubleSegment = 1<<segmentLength - 1

// The actual radix of a BinaryInteger
const carryThreshold DoubleSegment = 1 << segmentLength
const ten Segment = 10

const signPositive Segment = 0
const signNegative Segment = 1

// The maximum value of powers of 10 whose value is less than the maximum value of uint64
const maxIndex10 = "10000000000000000000"

// On average, the length of a number of radix 10 will be the multiplication of
// that of a number of this radix and this expandFactor
var expandFactor = math.Log(float64(carryThreshold)) / math.Log(float64(ten))
