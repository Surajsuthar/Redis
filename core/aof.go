// Append-only-file helpers for rewriting the current in-memory data set.
package core

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Surajsuthar/go-redis/config"
)

// dumpKey writes one key as a SET command into the AOF file.
func dumpKey(fp *os.File, key string, obj *Obj) {
	cmd := fmt.Sprintf("SET %s %s", key, obj.value)
	tokens := strings.Split(cmd, " ")
	fp.Write(Encode(tokens, false))
}

// dumpAllAOF rewrites all current keys into the configured AOF path.
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
