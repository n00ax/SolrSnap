package main

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"solr-snapshot-service/automated"
	"solr-snapshot-service/backup"
	"solr-snapshot-service/config"
	"strings"
)

func main() {
	fmt.Printf("%s CLI Mode, commit %s\n", config.ProductName, config.GitCommit)
	if len(os.Args) == 1 {
		backup.Help()
	} else {
		Cli()
	}
}

func Cli() {
	action := os.Args[1]
	if strings.EqualFold(action, "LoadRepoFromZK") && len(os.Args) == 6 {
		backup.LoadRepoFromZK(os.Args[2], os.Args[3], os.Args[4], os.Args[5])
	} else if strings.EqualFold(action, "LoadBlobFromZK") && len(os.Args) == 7 {
		backup.LoadBlobFromZK(os.Args[2], os.Args[3], os.Args[4], os.Args[5], os.Args[6])
	} else if strings.EqualFold(action, "LoadZKFromBlob") && len(os.Args) == 4 {
		backup.LoadZKFromBlob(os.Args[2], os.Args[3])
	} else if strings.EqualFold(action, "LoadRemoteRepoFromZK") && len(os.Args) == 6 {
		backup.LoadRemoteRepoFromZK(os.Args[2], os.Args[3], os.Args[4], os.Args[5])
	} else if strings.EqualFold(action, "BlobInfo") && len(os.Args) == 3 {
		backup.BlobInfo(os.Args[2])
	} else if strings.EqualFold(action, "CustomizeBlob") && len(os.Args) == 6 {
		backup.CustomizeBlob(os.Args[2], os.Args[3], os.Args[4], os.Args[5])
	} else if strings.EqualFold(action, "CreateCores") && len(os.Args) == 4 {
		backup.CreateCores(os.Args[2], os.Args[3])
	} else if strings.EqualFold(action, "Help") {
		backup.Help()
	} else if strings.EqualFold(action, "Automated") {
		automated.Cli()
	} else {
		fmt.Printf("Usage %s action [options]\n", os.Args[0])
		color.Red("Invalid command, enter \"Help\" for available commands and their usage")
		os.Exit(-1)
	}
	os.Exit(0)
}
