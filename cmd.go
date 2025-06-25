package main

import "os"

type CmdArgs struct {
	ConfigFile string
	DebugMode  bool
}

var Args = getArgs()

func getArgs() CmdArgs {
	args := os.Args
	if len(args) < 2 {
		panic("No config file provided")
	}
	return CmdArgs{
		ConfigFile: args[1],
		DebugMode:  Contains(args, "--debug"),
	}
}
