package automated

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"strings"
)

type command struct{
	description string
	argLength int
	runFunc func([]string)
}
var Commands = map[string]*command {
	"CreateLocal" : {
		description: "Creates Local Solr Instance From Remote Blob [blobPath]",
		argLength:   1,
		runFunc:     CreateLocal,
	},
}
func printAllArgs(){
	fmt.Printf("== Available job commands ==\n")
	for commandName, command := range Commands {
		fmt.Printf("%s : %s\n", commandName, command.description)
	}
}
func Cli() {
	color.Green("== Automated job mode ==")
	if len(os.Args) >= 3 {
		curCommand, exists := Commands[os.Args[2]]
		if !exists {
			for commandName, command := range Commands {
				if strings.EqualFold(commandName, os.Args[2]) {
					exists = true
					curCommand = command
				}
			}
		}
		if exists == false {
			color.Red("Invalid automated job command, %s", os.Args[2])
			printAllArgs()
			os.Exit(-1)
		}else if len(os.Args) - 3 != curCommand.argLength {
			color.Red("Invalid argument length for command %s, expected %d args, got %d args", os.Args[2], curCommand.argLength, len(os.Args)-3)
			color.Red("Description: %s\n", curCommand.description)
			os.Exit(-1)
		} else {
			curCommand.runFunc(os.Args[3:])
		}
	} else {
		printAllArgs()
	}
}