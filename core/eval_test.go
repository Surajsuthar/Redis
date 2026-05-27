package core

import (
	"bytes"
	"testing"
	"time"
)

func resetStore() {
	store = make(map[string]*Obj)
}

func TestEvalAndRespondSetGetAndDel(t *testing.T) {
	resetStore()

	cmds := RedisCmds{
		{Cmd: "SET", Args: []string{"name", "suraj"}},
		{Cmd: "GET", Args: []string{"name"}},
		{Cmd: "DEL", Args: []string{"name"}},
		{Cmd: "GET", Args: []string{"name"}},
	}

	var conn bytes.Buffer
	EvalAndRespond(&cmds, &conn)

	const want = "+Ok\r\n$5\r\nsuraj\r\n:1\r\n$-1\r\n"
	if conn.String() != want {
		t.Fatalf("unexpected response:\nwant %q\n got %q", want, conn.String())
	}
}

func TestSetWithEXExpiresKey(t *testing.T) {
	resetStore()

	if got := string(evalSet([]string{"session", "token", "EX", "1"}, nil)); got != "+Ok\r\n" {
		t.Fatalf("unexpected SET response: %q", got)
	}

	if got := string(evalGet([]string{"session"}, nil)); got != "$5\r\ntoken\r\n" {
		t.Fatalf("expected key before expiry, got %q", got)
	}

	store["session"].ExpireAt = time.Now().UnixMilli() - 1

	if got := string(evalGet([]string{"session"}, nil)); got != "$-1\r\n" {
		t.Fatalf("expected missing key after expiry, got %q", got)
	}
}

func TestExpireAndTTL(t *testing.T) {
	resetStore()

	PUT("name", NewObj("suraj", -1))

	if got := string(evalEXPIRE([]string{"name", "1"}, nil)); got != ":1\r\n" {
		t.Fatalf("unexpected EXPIRE response: %q", got)
	}

	ttl := string(evalTTL([]string{"name"}, nil))
	if ttl != ":0\r\n" && ttl != ":1\r\n" {
		t.Fatalf("expected TTL to be 0 or 1 second, got %q", ttl)
	}
}
