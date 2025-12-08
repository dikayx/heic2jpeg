package app

import "os"

func Run() {
	if len(os.Args) == 1 {
		runGuided()
		return
	}

	switch os.Args[1] {
	case "cli":
		runCLI(os.Args[2:])
	case "help", "--help", "-h":
		printHelp()
	default:
		runGuided()
	}
}
