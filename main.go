package main

import (
	"fmt"
	"os"

	"github.com/LilRooks/ytdl-ipfs-archiver/internal/app"
)

func main() {
	if err := app.Run(os.Args, os.Stdout); err.Error != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error)
		os.Exit(err.Code)
	}
}
