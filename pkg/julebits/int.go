package julebits

import (
	"math"
	"strconv"
	"strings"
)

// MAX_INT is the maximum bitsize of integer types.
const MAX_INT = 64

type bitChecker = func(val string, base, bit int) error

// CheckBitInt reports integer is compatible this bit-size or not.
func CheckBitInt(val string, bit int) bool {
	return checkBit(val, bit, func(val string, base, bit int) error {
		_, err := strconv.ParseInt(val, base, bit)
		return err
	})
}

// CheckBitUInt reports unsigned integer is
// compatible this bit-size or not.
func CheckBitUInt(val string, bit int) bool {
	return checkBit(val, bit, func(val string, base, bit int) error {
		_, err := strconv.ParseUint(val, base, bit)
		return err
	})
}

func checkBit(val string, bit int, checker bitChecker) bool {
	var err error
	switch {
	case val == "":
		return false
	case len(val) == 1:
		return true
	case strings.HasPrefix(val, "0x"):
		err = checker(val[2:], 16, bit)
	case strings.HasPrefix(val, "0b"):
		err = checker(val[2:], 2, bit)
	case val[0] == '0':
		err = checker(val[1:], 8, bit)
	default:
		err = checker(val, 10, bit)
	}
	return err == nil
}

// BitsizeInt returns minimum bitsize of given value.
func BitsizeInt(x int64) uint64 {
	switch {
	case x >= math.MinInt8 && x <= math.MaxInt8:
		return 8
	case x >= math.MinInt16 && x <= math.MaxInt16:
		return 16
	case x >= math.MinInt32 && x <= math.MaxInt32:
		return 32
	default:
		return MAX_INT
	}
}

// BitsizeUInt returns minimum bitsize of given value.
func BitsizeUInt(x uint64) uint64 {
	switch {
	case x <= math.MaxUint8:
		return 8
	case x <= math.MaxUint16:
		return 16
	case x <= math.MaxUint32:
		return 32
	default:
		return MAX_INT
	}
}
