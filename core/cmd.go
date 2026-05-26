package core

type RedisCLI struct {
	Cmd  string
	Args []string
}

type RedisCmds []RedisCLI
