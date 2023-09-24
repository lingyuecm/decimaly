package model

import "math"

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
