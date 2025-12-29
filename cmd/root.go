package cmd

import (
	"flag"
	"fmt"

	"github.com/vanshitkumar/sappend/internal"
)

func RootCommand() error {
	fileUrl := flag.String("file-url", "", "URL of the file to synchronize")
	localPath := flag.String("local-file", "", "Local path to synchronize the file to")

	flag.Parse()

	if *fileUrl == "" {
		return fmt.Errorf("Incorrect usage: file-url flag is required but not provided")
	}

	return internal.SyncFile(*fileUrl, *localPath)
}
