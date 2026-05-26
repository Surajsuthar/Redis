package core

import (
	"errors"
	"fmt"
)

func readSimpleString(data []byte) (string, int, error) {
	pos := 1
	for ; data[pos] != '\r'; pos++ {
	}
	return string(data[1:pos]), pos + 2, nil
}

func readError(data []byte) (string, int, error) {
	return readSimpleString(data)
}

func readInt(data []byte) (int, int, error) {
	pos := 1
	var value int = 0
	for ; data[pos] != '\r'; pos++ {
		value = value*10 + int(data[pos]-'0')
	}
	return value, pos + 2, nil
}

func readLen(data []byte) (int, int) {
	pos, len := 0, 0
	for pos = range data {
		if !(data[pos] >= '0' && data[pos] <= '9') {
			return len, pos + 2
		}
		len = len*10 + int(data[pos]-'0')
	}
	return 0, 0
}

func readBulkString(data []byte) (string, int, error) {
	pos := 1
	len, gap := readLen(data[pos:])
	pos += gap
	return string(data[pos:(pos + len)]), pos + len + 2, nil
}

func readArray(data []byte) (interface{}, int, error) {
	pos := 1
	len, delta := readLen(data[pos:])
	pos += delta

	var elems []interface{} = make([]interface{}, len)

	for i := range elems {
		ele, delta, err := DecodeOne(data[pos:])
		if err != nil {
			return nil, 0, err
		}
		elems[i] = ele
		pos += delta
	}
	fmt.Println("elems", elems, "pos", pos, "delta", delta)
	return elems, pos, nil
}

func DecodeOne(data []byte) (interface{}, int, error) {
	if len(data) == 0 {
		return nil, 0, errors.New("No data")
	}
	switch data[0] {
	case '+':
		return readSimpleString(data)
	case ':':
		return readInt(data)
	case '-':
		return readError(data)
	case '$':
		return readBulkString(data)
	case '*':
		return readArray(data)
	}
	return nil, 0, nil
}

func Decode(data []byte) ([]interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("No data")
	}

	var value []interface{} = make([]interface{}, 0)
	var index int = 0

	for index < len(data) {
		val, gap, err := DecodeOne(data[index:])
		if err != nil {
			return value, err
		}
		index = index + gap
		value = append(value, val)
	}

	return value, nil
}

// func DecodeArrayString(data []byte) ([]string, error) {s
// 	value, err := Decode(data)
// 	if err != nil {
// 		return nil, err
// 	}

// 	ts := value.([]interface{})
// 	token := make([]string, len(ts))
// 	for i := range token {
// 		token[i] = ts[i].(string)
// 	}

// 	return token, nil
// }

func Encode(value interface{}, isSimple bool) []byte {
	switch v := value.(type) {
	case string:
		if isSimple {
			return []byte(fmt.Sprintf("+%s\r\n", v))
		}
		return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(v), v))
	case int, int8, int16, int32, int64:
		return []byte(fmt.Sprintf(":%d\r\n", v))
	case error:
		return []byte(fmt.Sprintf("-%s\r\n", v))
	case nil:
		return []byte("$-1\r\n")
	}
	return []byte{}
}
