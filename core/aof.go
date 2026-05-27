package core

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Surajsuthar/go-redis/config"
)

func dumpKey(fp *os.File, key string, obj *Obj) {
	cmd := fmt.Sprintf("SET %s %s", key, obj.value)
	tokens := strings.Split(cmd, " ")
	fp.Write(Encode(tokens, false))
}

func dumpAllAOF() {
	fp, err := os.OpenFile(config.AofPath, os.O_CREATE|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Print("error", err)
		return
	}
	log.Panicln("rewriting to File")
	for key, obj := range store {
		dumpKey(fp, key, obj)
	}
}
