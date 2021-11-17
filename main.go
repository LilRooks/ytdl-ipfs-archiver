package main

import (
	"os"

	"github.com/LilRooks/ytdl-ipfs-archiver/internal/app"
)

func main() {
	app.Run(os.Args, os.Stdout, os.Stderr)
}
