// Type and encoding helpers for compact object metadata.
package core

import (
	"errors"
	"strconv"
)

// getType extracts the object type bits from the combined type/encoding byte.
func getType(te uint8) uint8 {
	return (te >> 4) << 4
}

// getEncoding extracts the object encoding bits from the combined byte.
func getEncoding(te uint8) uint8 {
	return te & 0b00001111
}

// assertType rejects operations that are not valid for the object type.
func assertType(te uint8, t uint8) error {
	if getType(te) != t {
		return errors.New("not permited on this type")
	}
	return nil
}

// assertEocoding rejects operations that require a different object encoding.
func assertEocoding(te uint8, e uint8) error {
	if getEncoding(te) != e {
		return errors.New("not permited on this type")
	}
	return nil
}

// deduceTyeEncoding chooses the string object encoding from the value shape.
func deduceTyeEncoding(v string) (uint8, uint8) {
	oType := OBJ_TYPE_STRING
	if _, err := strconv.ParseInt(v, 10, 64); err == nil {
		return oType, OBJ_ENCODING_INT
	}
	if len(v) <= 44 {
		return oType, OBJ_ENCODING_EMBSTR
	}
	return oType, OBJ_ENCODING_RAW
}
