package xtype

import "github.com/the-xlang/xxc/lex/tokens"

// Data type (built-in) constants.
const (
	Void    uint8 = 0
	I8      uint8 = 1
	I16     uint8 = 2
	I32     uint8 = 3
	I64     uint8 = 4
	U8      uint8 = 5
	U16     uint8 = 6
	U32     uint8 = 7
	U64     uint8 = 8
	Bool    uint8 = 9
	Str     uint8 = 10
	F32     uint8 = 11
	F64     uint8 = 12
	Any     uint8 = 13
	Char    uint8 = 14
	Id      uint8 = 15
	Func    uint8 = 16
	Nil     uint8 = 17
	Size    uint8 = 18
	Map     uint8 = 19
	Voidptr uint8 = 20
	Enum    uint8 = 21
)

const (
	NumericTypeStr = "<numeric>"
	NilTypeStr     = "<nil>"
	VoidTypeStr    = "<void>"
)

// TypeGreaterThan reports type one is greater than type two or not.
func TypeGreaterThan(t1, t2 uint8) bool {
	switch t1 {
	case I16:
		return t2 == I8
	case I32:
		return t2 == I8 ||
			t2 == I16
	case I64:
		return t2 == I8 ||
			t2 == I16 ||
			t2 == I32
	case U16:
		return t2 == U8
	case U32:
		return t2 == U8 ||
			t2 == U16
	case U64:
		return t2 == U8 ||
			t2 == U16 ||
			t2 == U32
	case F32:
		return t2 != Any && t2 != F64
	case F64:
		return t2 != Any
	case Size, Enum:
		return true
	}
	return false
}

// TypeAreCompatible reports type one and type two is compatible or not.
func TypesAreCompatible(t1, t2 uint8, ignoreany bool) bool {
	if !ignoreany && t2 == Any {
		return true
	}
	switch t1 {
	case Any:
		if ignoreany {
			return false
		}
		return true
	case I8:
		return t2 == I8
	case I16:
		return t2 == I8 ||
			t2 == I16 ||
			t2 == Char
	case I32:
		return t2 == I8 ||
			t2 == I16 ||
			t2 == I32 ||
			t2 == I64 ||
			t2 == F32 ||
			t2 == F64 ||
			t2 == Size ||
			t2 == Char
	case I64:
		return t2 == I8 ||
			t2 == I16 ||
			t2 == I32 ||
			t2 == I64 ||
			t2 == F32 ||
			t2 == F64 ||
			t2 == Size ||
			t2 == Char
	case U8:
		return t2 == U8 ||
			t2 == U16 ||
			t2 == U32 ||
			t2 == U64 ||
			t2 == F32 ||
			t2 == F64 ||
			t2 == Size ||
			t2 == Char
	case U16:
		return t2 == U16 ||
			t2 == U32 ||
			t2 == U64 ||
			t2 == F32 ||
			t2 == F64 ||
			t2 == Size ||
			t2 == Char
	case U32:
		return t2 == U32 ||
			t2 == U64 ||
			t2 == F32 ||
			t2 == F64 ||
			t2 == Size ||
			t2 == Char
	case U64, Size:
		return t2 == U64 ||
			t2 == F32 ||
			t2 == F64 ||
			t2 == Size ||
			t2 == Char
	case Bool:
		return t2 == Bool
	case Str:
		return t2 == Str
	case F32:
		return t2 == F32 ||
			t2 == I8 ||
			t2 == I16 ||
			t2 == I32 ||
			t2 == U8 ||
			t2 == U16 ||
			t2 == U32 ||
			t2 == Char
	case F64:
		return t2 == F64 ||
			t2 == F32 ||
			t2 == I8 ||
			t2 == I16 ||
			t2 == I32 ||
			t2 == U8 ||
			t2 == U16 ||
			t2 == U32 ||
			t2 == Char
	case Char:
		return t2 == Char ||
			t2 == U8
	case Nil:
		return t2 == Nil
	}
	return false
}

// IsIntegerType reports type is signed/unsigned integer or not.
func IsIntegerType(t uint8) bool {
	return IsSignedIntegerType(t) || IsUnsignedNumericType(t)
}

// IsNumericType reports type is numeric or not.
func IsNumericType(t uint8) bool {
	return IsIntegerType(t) || IsFloatType(t)
}

// IsFloatType reports type is float or not.
func IsFloatType(t uint8) bool {
	return t == F32 || t == F64
}

// IsSignedNumericType reports type is signed numeric or not.
func IsSignedNumericType(t uint8) bool {
	return IsSignedIntegerType(t) ||
		t == F32 ||
		t == F64
}

// IsSignedIntegerType reports type is signed integer or not.
func IsSignedIntegerType(t uint8) bool {
	return t == I8 ||
		t == I16 ||
		t == I32 ||
		t == I64
}

// IsUnsignedNumericType reports type is unsigned numeric or not.
func IsUnsignedNumericType(t uint8) bool {
	return t == U8 ||
		t == U16 ||
		t == U32 ||
		t == U64 ||
		t == Size
}

// TypeFromId returns type id of specified type code.
func TypeFromId(id string) uint8 {
	switch id {
	case tokens.I8, tokens.SBYTE:
		return I8
	case tokens.I16:
		return I16
	case tokens.I32:
		return I32
	case tokens.I64:
		return I64
	case tokens.U8, tokens.BYTE:
		return U8
	case tokens.U16:
		return U16
	case tokens.U32:
		return U32
	case tokens.U64:
		return U64
	case tokens.STR:
		return Str
	case tokens.BOOL:
		return Bool
	case tokens.F32:
		return F32
	case tokens.F64:
		return F64
	case "any":
		return Any
	case tokens.CHAR:
		return Char
	case tokens.SIZE:
		return Size
	case tokens.VOIDPTR:
		return Voidptr
	}
	return 0
}

func CxxTypeIdFromType(typeCode uint8) string {
	switch typeCode {
	case Void:
		return "void"
	case I8:
		return tokens.I8
	case I16:
		return tokens.I16
	case I32:
		return tokens.I32
	case I64:
		return tokens.I64
	case U8:
		return tokens.U8
	case U16:
		return tokens.U16
	case U32:
		return tokens.U32
	case U64:
		return tokens.U64
	case Bool:
		return tokens.BOOL
	case F32:
		return tokens.F32
	case F64:
		return tokens.F64
	case Any:
		return "any"
	case Str:
		return tokens.STR
	case Char:
		return tokens.CHAR
	case Size:
		return tokens.SIZE
	case Voidptr:
		return tokens.VOIDPTR
	}
	return "" // Unreachable code.
}

// DefaultValOfType returns default value of specified type.
//
// Special case is:
//  DefaultValOfType(t) = "nil" if t is invalid
//  DefaultValOfType(t) = "nil" if t is not have default value
func DefaultValOfType(code uint8) string {
	if IsNumericType(code) {
		return "0"
	}
	switch code {
	case Bool:
		return "false"
	case Str:
		return `""`
	case Char:
		return `'\0'`
	}
	return "nil"
}