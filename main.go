package main

import (
	"fmt"
	"os"

	"github.com/LilRooks/ytdl-ipfs-archiver/internal/app"
)

func main() {
	code, err := app.Run(os.Args, os.Stdout, os.Stderr)
	if err != nil {
		fmt.Fprintf(os.Stdout, "%s", err)
		os.Exit(code)
	}
}
