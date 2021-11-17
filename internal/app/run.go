package app

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"io"
	"io/fs"
	"io/ioutil"

	"github.com/LilRooks/ytdl-ipfs-archiver/internal/pkg/config"
	"github.com/LilRooks/ytdl-ipfs-archiver/internal/pkg/ipfs"
	"github.com/LilRooks/ytdl-ipfs-archiver/internal/pkg/table"
	"github.com/LilRooks/ytdl-ipfs-archiver/internal/pkg/ytdl"

	"github.com/ipfs/go-cid"
	"github.com/web3-storage/go-w3s-client"
)

var (
	ytdlPath string
	confPath string
	tablPath string
)

const (
	errorGeneric = iota
	errorConfig
	errorYTDL
	errorTable
	errorIPFS
)

func handleErr(logger *log.Logger, stderr io.Writer, err error, exitCode int) {
	if err != nil {
		fmt.Fprintf(stderr, "%s\n", err)
		os.Exit(exitCode)
	}
	logger.SetPrefix("[base] ")
}

// Run is the actual code for the command
func Run(args []string, stdout io.Writer, stderr io.Writer) {
	logger := log.New(stdout, "[base] ", log.LstdFlags)
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)

	flags.StringVar(&ytdlPath, "bin", "", "path to the youtube-dl binary (defaults to one in PATH)")
	flags.StringVar(&confPath, "cfg", "", "path to the configuration file to use")
	flags.StringVar(&tablPath, "tab", "./table.sqlite", "path to the table file to use")

	//if err := flags.Parse(args[1:]); err != nil {
	//	return err, errorGeneric
	//}

	err := flags.Parse(args[1:])
	handleErr(logger, stderr, err, errorConfig)
	ytdlOptions := flags.Args()

	// Read configuration into configs variable
	err, configs := config.Parse(confPath)
	handleErr(logger, stderr, err, errorConfig)

	// Configuration file is the real default
	if len(ytdlPath) == 0 {
		ytdlPath = configs.Binary
	}
	removeYTDL := (ytdlPath == "embedded")

	// Check binary is there
	err, ytdlPath = ytdl.ParsePath(ytdlPath)
	handleErr(logger, stderr, err, errorYTDL)

	if removeYTDL {
		defer os.Remove(ytdlPath)
	}
	// Get the keys needed for the table
	var (
		filename string
		id       string
		format   string
		location string
	)
	_, errNotExist := os.Stat(tablPath)

	err, db := table.OpenDB(tablPath)
	handleErr(logger, stderr, err, errorTable)
	defer db.Close()

	// Only initialized if the file did not originally exist
	if errors.Is(errNotExist, os.ErrNotExist) {
		errInit := table.InitializeTable(db)
		handleErr(logger, stderr, errInit, errorTable)
	}

	err, id, format = ytdl.GetIdentifiers(logger, ytdlPath, ytdlOptions)
	handleErr(logger, stderr, err, errorYTDL)

	err, location = table.Fetch(db, id, format)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		handleErr(logger, stderr, err, errorTable)
	}

	cid, _ := cid.Decode(location)

	if val, ok := os.LookupEnv("TOKEN"); ok {
		configs.Token = val
	}
	c, err := w3s.NewClient(w3s.WithToken(configs.Token))
	handleErr(logger, stderr, err, errorConfig)

	if len(location) == 0 {
		err := ytdl.Download(logger, ytdlPath, ytdlOptions)
		handleErr(logger, stderr, err, errorYTDL)

		err, filename = ytdl.GetFilename(logger, ytdlPath, ytdlOptions)
		handleErr(logger, stderr, err, errorYTDL)

		// only uploads first match, may have undefined behavior if file of same name exists
		// this is why file is stored with identifying information
		filenames, _ := filepath.Glob(filename + ".*")
		filename = filenames[0]

		err, location = ipfs.Store(c, filename)
		handleErr(logger, stderr, err, errorIPFS)

		err = table.Store(db, id, format, location)
		handleErr(logger, stderr, err, errorTable)
	} else {
		logger.Printf("[w3s] Getting %s\n", location)
		res, _ := c.Get(context.Background(), cid)

		// res is a http.Response with an extra method for reading IPFS UnixFS files!
		f, fsys, _ := res.Files()

		// Download directory contents
		if d, ok := f.(fs.ReadDirFile); ok {
			ents, _ := d.ReadDir(0)
			for _, ent := range ents {
				filename = ent.Name()
				file, err := fsys.Open("/" + ent.Name())
				handleErr(logger, stderr, err, errorIPFS)

				data, err := ioutil.ReadAll(file)
				handleErr(logger, stderr, err, errorIPFS)

				err = ioutil.WriteFile(ent.Name(), data, 0755)
				handleErr(logger, stderr, err, errorIPFS)
			}
		}
	}

	logger.Printf("File is available locally at %s\n", filename)
	logger.Printf("File is also available at https://%s.ipfs.dweb.link\n", location)
}
