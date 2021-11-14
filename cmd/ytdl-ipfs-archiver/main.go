package main

import (
	"flag"
	"fmt"
)

var ytdlBinary string
var configPath string
var tablePath string

func init() {
	flag.StringVar(&ytdlBinary, "bin", "", "path to the youtube-dl binary (defaults to one in PATH)")
	flag.StringVar(&configPath, "cfg", "", "path to the configuration file to use")
	flag.StringVar(&tablePath, "tab", "./table.edn", "path to the table file to use")
}

func main() {
	flag.Parse()
	ytdlOptions := flag.Args()
	for _, v := range ytdlOptions {
		fmt.Print(v + " ")
	}
}
