package main

import (
	"fmt"
	"os"

	"github.com/LilRooks/ytdl-ipfs-archiver/internal/cmd"
)

func main() {
	if err := cmd.Run(os.Args, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error)
		return err.Code
	}
}
