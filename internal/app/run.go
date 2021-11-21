package app

import (
	"errors"
	"flag"
	"log"
	"os"
	"path/filepath"

	"io"

	"github.com/LilRooks/ytdl-ipfs-archiver/internal/pkg/config"
	"github.com/LilRooks/ytdl-ipfs-archiver/internal/pkg/ipfs"
	"github.com/LilRooks/ytdl-ipfs-archiver/internal/pkg/table"
	"github.com/LilRooks/ytdl-ipfs-archiver/internal/pkg/ytdl"

	"github.com/web3-storage/go-w3s-client"
)

var (
	ytdlPath string
	confPath string
	tablPath string
	pinAddDB bool
	remoteDB bool
)

const (
	errorNone = iota
	errorConfig
	errorYTDL
	errorTable
	errorIPFS
)

// Run is the actual code for the command
func Run(args []string, stdout io.Writer, stderr io.Writer) (int, error) {
	logger := log.New(stdout, "[base] ", log.Ltime|log.Lmsgprefix)
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)

	flags.StringVar(&ytdlPath, "bin", "", "path to the youtube-dl binary (defaults to one in PATH)")
	flags.StringVar(&confPath, "cfg", "", "path to the configuration file to use")
	flags.StringVar(&tablPath, "tab", "./table.sqlite", "path to the table file to use")
	flags.BoolVar(&pinAddDB, "pin", false, "pin the database to the local daemon")
	flags.BoolVar(&remoteDB, "rdb", false, "pulls remote database to `pwd`")

	err := flags.Parse(args[1:])
	if err != nil {
		return errorConfig, err
	}
	ytdlOptions := flags.Args()

	// Read configuration into configs variable
	configs, err := config.Parse(confPath)
	if err != nil {
		return errorConfig, err
	}

	// Configuration file is the real default
	if len(ytdlPath) == 0 {
		ytdlPath = configs.Binary
	}

	// Check binary is there
	ytdlPath, err = ytdl.ParsePath(ytdlPath)
	if err != nil {
		return errorYTDL, err
	}

	// Get the keys needed for the table
	var (
		filename string
		id       string
		format   string
		location string
	)
	c, err := w3s.NewClient(w3s.WithToken(configs.Token))
	if err != nil {
		return errorConfig, err
	}
	tablPath, err = ipfs.Fetch(c, tablPath)
	if err != nil {
		return errorIPFS, err
	}
	_, errTableExist := os.Stat(tablPath)

	db, err := table.OpenDB(tablPath)
	if err != nil {
		return errorTable, err
	}
	defer db.Close()

	// Only initialized if the file did not originally exist and is local only
	if errors.Is(errTableExist, os.ErrNotExist) {
		err := table.InitializeTable(db)
		if err != nil {
			return errorTable, err
		}

	}

	id, format, err = ytdl.GetIdentifiers(logger, ytdlPath, ytdlOptions)
	if err != nil {
		return errorYTDL, err
	}

	location, err = table.Fetch(db, id, format)
	if err != nil {
		return errorTable, err
	}

	if val, ok := os.LookupEnv("TOKEN"); ok {
		configs.Token = val
	}

	if len(location) == 0 {
		err := ytdl.Download(logger, ytdlPath, ytdlOptions)
		if err != nil {
			return errorYTDL, err
		}

		filename, err = ytdl.GetFilename(logger, ytdlPath, ytdlOptions)
		if err != nil {
			return errorYTDL, err
		}

		// only uploads first match, may have undefined behavior if file of same name exists
		// this is why file is stored with identifying information
		filenames, _ := filepath.Glob(filename + ".*")
		filename = filenames[0]

		location, err = ipfs.Store(c, filename)
		if err != nil {
			return errorIPFS, err
		}

		err = table.Store(db, id, format, location)
		if err != nil {
			return errorIPFS, err
		}
	} else {
		var err error
		filename, err = ipfs.Fetch(c, location)
		if err != nil {
			return errorIPFS, err
		}
	}

	logger.SetPrefix("[base] ")
	logger.Printf("File is available locally at %s\n", filename)
	logger.Printf("File is also available at https://%s.ipfs.dweb.link\n", location)
	return errorNone, nil
}
