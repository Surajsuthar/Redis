package core

import (
	"errors"
	"io"
	"strconv"
	"time"
)

func evalPing(args []string, conn io.ReadWriter) error {
	if len(args) >= 2 {
		return errors.New("Err wrong number of cammand for ping cammand")
	}
	var b []byte
	if len(args) == 0 {
		b = Encode("PONG", true)
	} else {
		b = Encode(args[0], false)
	}

	_, err := conn.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func evalSet(args []string, conn io.ReadWriter) error {
	if len(args) <= 1 {
		return errors.New("(Error) wroung number of argument to set")
	}

	var key, value string
	var onDurationOnMs int64 = -1

	key, value = args[0], args[1]

	for i := 2; i < len(args); i++ {
		switch args[i] {
		case "EX", "ex", "Ex":
			i++
			if i == len(args) {
				return errors.New("(Error) syntax error")
			}
			onDurationOnSec, err := strconv.ParseInt(args[3], 10, 64)
			if err != nil {
				return errors.New("(Error) value is not integer or out of range")
			}

			onDurationOnMs = onDurationOnSec * 100

		default:
			return errors.New("(Error) syntax error")
		}
	}

	PUT(key, NewObj(value, onDurationOnMs))
	conn.Write([]byte("+Ok\r\n"))
	return nil
}

func evalGet(args []string, conn io.ReadWriter) error {
	if len(args) != 1 {
		return errors.New("(Error) wroung number of argument to set")
	}

	var key string = args[0]
	obj := GET(key)

	if obj == nil {
		conn.Write([]byte("$-1\r\n"))
		return nil
	}

	if obj.ExpireAt != -1 || obj.ExpireAt <= time.Now().Local().UnixMilli() {
		conn.Write([]byte("$-1\r\n"))
		return nil
	}

	conn.Write(Encode(obj.Value, false))
	return nil
}

func evalTTL(args []string, conn io.ReadWriter) error {
	if len(args) != 1 {
		return errors.New("(Error) wroung number of argument to set")
	}

	var key string = args[0]
	obj := GET(key)

	if obj == nil {
		conn.Write([]byte(":-2\r\n"))
		return nil
	}

	if obj.ExpireAt == -1 {
		conn.Write([]byte(":-1\r\n"))
		return nil
	}

	duration := obj.ExpireAt - time.Now().UnixMilli()

	if duration < 0 {
		conn.Write([]byte(":-2\r\n"))
		return nil
	}

	conn.Write(Encode(int(duration/1000), false))
	return nil
}

func evalDEL(args []string, conn io.ReadWriter) error {
	var countdelete int = 0

	for _, key := range args {
		if ok := DEL(key); ok {
			countdelete++
		}
	}

	conn.Write(Encode(countdelete, false))
	return nil
}

func evalEXPIRE(args []string, conn io.ReadWriter) error {
	if len(args) <= 1 {
		return errors.New("(Error) wroung number of argument to set")
	}

	var key string = args[0]
	exDuruation, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return errors.New("(Error) value is not integer or out of range")
	}

	obj := GET(key)

	if obj == nil {
		conn.Write([]byte("$0\r\n"))
		return nil
	}
	obj.ExpireAt = time.Now().UnixMilli() + exDuruation*100

	conn.Write([]byte(":1\r\n"))
	return nil
}

func EvalAndRespond(cmd *RedisCLI, conn io.ReadWriter) error {
	switch cmd.Cmd {
	case "PING":
		return evalPing(cmd.Args, conn)
	case "SET":
		return evalSet(cmd.Args, conn)
	case "GET":
		return evalGet(cmd.Args, conn)
	case "TTL":
		return evalTTL(cmd.Args, conn)
	case "DEL":
		return evalDEL(cmd.Args, conn)
	case "EXPIRE":
		return evalEXPIRE(cmd.Args, conn)
	default:
		return evalPing(cmd.Args, conn)
	}
}
