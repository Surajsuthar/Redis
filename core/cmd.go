// Parsed command types shared between the TCP and evaluator layers.
package core

type RedisCLI struct {
	Cmd  string
	Args []string
}

type RedisCmds []RedisCLI
