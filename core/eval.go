// Command evaluation for Redis-like commands supported by this server.
package core

import (
	"bytes"
	"errors"
	"io"
	"strconv"
	"time"
)

// evalPing handles PING with optional echo payload.
func evalPing(args []string, conn io.ReadWriter) []byte {
	if len(args) >= 2 {
		return Encode(errors.New("Err wrong number of cammand for ping cammand"), false)
	}
	var b []byte
	if len(args) == 0 {
		b = Encode("PONG", true)
	} else {
		b = Encode(args[0], false)
	}

	return b
}

// evalSet stores a string value and optionally attaches an EX expiry.
func evalSet(args []string, conn io.ReadWriter) []byte {
	if len(args) <= 1 {
		return Encode(errors.New("(Error) wroung number of argument to set"), false)
	}

	var key, value string
	var onDurationOnMs int64 = -1

	key, value = args[0], args[1]
	te, e := deduceTyeEncoding(value)

	for i := 2; i < len(args); i++ {
		switch args[i] {
		case "EX", "ex", "Ex":
			i++
			if i == len(args) {
				return Encode(errors.New("(Error) syntax error"), false)
			}
			onDurationOnSec, err := strconv.ParseInt(args[i], 10, 64)
			if err != nil {
				return Encode(errors.New("(Error) value is not integer or out of range"), false)
			}

			onDurationOnMs = onDurationOnSec * 1000

		default:
			return Encode(errors.New("(Error) syntax error"), false)
		}
	}

	PUT(key, NewObj(value, onDurationOnMs, te, e))
	return []byte("+Ok\r\n")
}

// evalGet returns the stored value for a key or a RESP nil bulk string.
func evalGet(args []string, conn io.ReadWriter) []byte {
	if len(args) != 1 {
		return Encode(errors.New("(Error) wroung number of argument to set"), false)
	}

	var key string = args[0]
	obj := GET(key)

	if obj == nil {
		return []byte("$-1\r\n")
	}

	if obj.ExpireAt != -1 && obj.ExpireAt <= time.Now().Local().UnixMilli() {
		return []byte("$-1\r\n")
	}

	return Encode(obj.value, false)
}

// evalTTL returns the remaining lifetime of a key in whole seconds.
func evalTTL(args []string, conn io.ReadWriter) []byte {
	if len(args) != 1 {
		return Encode(errors.New("(Error) wroung number of argument to set"), false)
	}

	var key string = args[0]
	obj := GET(key)

	if obj == nil {
		return []byte(":-2\r\n")
	}

	if obj.ExpireAt == -1 {
		return []byte(":-1\r\n")
	}

	duration := obj.ExpireAt - time.Now().UnixMilli()

	if duration < 0 {
		return []byte(":-2\r\n")
	}

	return Encode(int(duration/1000), false)
}

// evalDEL deletes all requested keys and returns the number removed.
func evalDEL(args []string, conn io.ReadWriter) []byte {
	var countdelete int = 0

	for _, key := range args {
		if ok := DEL(key); ok {
			countdelete++
		}
	}

	return Encode(countdelete, false)
}

// evalEXPIRE attaches a new expiry time to an existing key.
func evalEXPIRE(args []string, conn io.ReadWriter) []byte {
	if len(args) <= 1 {
		return Encode(errors.New("(Error) wroung number of argument to set"), false)
	}

	var key string = args[0]
	exDuruation, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return Encode(errors.New("(Error) value is not integer or out of range"), false)
	}

	obj := GET(key)

	if obj == nil {
		return []byte("$0\r\n")
	}
	obj.ExpireAt = time.Now().UnixMilli() + exDuruation*1000

	return []byte(":1\r\n")
}

// evalBGREWRITEAOF rewrites the append-only file from the current store.
func evalBGREWRITEAOF(args []string, conn io.ReadWriter) []byte {
	dumpAllAOF()
	return []byte(":1\r\n")
}

// evalINCR increments an integer-encoded string value.
func evalINCR(args []string, conn io.ReadWriter) []byte {
	if len(args) != 1 {
		return Encode(errors.New("(Error) wrong number of arguments"), false)
	}

	var key string = args[0]
	obj := GET(key)

	if obj == nil {
		obj = NewObj("0", -1, OBJ_TYPE_STRING, OBJ_ENCODING_INT)
		PUT(key, obj)
	}

	if err := assertType(obj.TypeEncoding, OBJ_TYPE_STRING); err != nil {
		return Encode(err, false)
	}

	if err := assertEocoding(obj.TypeEncoding, OBJ_ENCODING_INT); err != nil {
		return Encode(err, false)
	}
	i, _ := strconv.ParseInt(obj.value.(string), 10, 64)
	i++
	obj.value = strconv.FormatInt(i, 10)

	return Encode(i, false)
}

// EvalAndRespond evaluates a batch of commands and writes one combined response.
func EvalAndRespond(cmd *RedisCmds, conn io.ReadWriter) {
	var response []byte
	buf := bytes.NewBuffer(response)

	for _, c := range *cmd {
		switch c.Cmd {
		case "PING":
			buf.Write(evalPing(c.Args, conn))
		case "SET":
			buf.Write(evalSet(c.Args, conn))
		case "GET":
			buf.Write(evalGet(c.Args, conn))
		case "TTL":
			buf.Write(evalTTL(c.Args, conn))
		case "DEL":
			buf.Write(evalDEL(c.Args, conn))
		case "EXPIRE":
			buf.Write(evalEXPIRE(c.Args, conn))
		case "BGREWRITEAOF":
			buf.Write(evalBGREWRITEAOF(c.Args, conn))
		case "INCR":
			buf.Write(evalINCR(c.Args, conn))
		default:
			buf.Write(evalPing(c.Args, conn))
		}
	}
	conn.Write(buf.Bytes())
}
