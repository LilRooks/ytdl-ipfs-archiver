package app

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"io"
	"io/fs"
	"io/ioutil"

	"github.com/LilRooks/ytdl-ipfs-archiver/internal/pkg/config"
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

// Run is the actual code for the command
func Run(args []string, stdout io.Writer) (error, int) {
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)

	flags.StringVar(&ytdlPath, "bin", "", "path to the youtube-dl binary (defaults to one in PATH)")
	flags.StringVar(&confPath, "cfg", "", "path to the configuration file to use")
	flags.StringVar(&tablPath, "tab", "./table.sqlite", "path to the table file to use")

	if err := flags.Parse(args[1:]); err != nil {
		return err, errorGeneric
	}
	ytdlOptions := flags.Args()

	// Read configuration into configs variable
	err, configs := config.Parse(confPath)
	if err != nil {
		return err, errorConfig
	}

	// Configuration file is the real default
	if len(ytdlPath) == 0 {
		ytdlPath = configs.Binary
	}
	removeYTDL := (ytdlPath == "embedded")

	// Check binary is there
	err, ytdlPath = ytdl.ParsePath(ytdlPath)
	if err != nil {
		return err, errorYTDL
	}

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
	err, id, format = ytdl.GetIdentifiers(ytdlPath, ytdlOptions)
	if err != nil {
		return err, errorYTDL
	}
	err, filename = ytdl.GetFilename(ytdlPath, ytdlOptions)
	if err != nil {
		return err, errorYTDL
	}

	_, errNotExist := os.Stat(tablPath)

	err, db := table.OpenDB(tablPath)
	if err != nil {
		return err, errorTable
	}
	defer db.Close()

	// Only initialized if the file did not originally exist
	if errors.Is(errNotExist, os.ErrNotExist) {
		errInit := table.InitializeTable(db)
		if errInit != nil {
			return errInit, errorTable
		}
	}

	err, location = table.Fetch(db, id, format)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err, errorTable
	}

	cid, _ := cid.Decode(location)

	if val, ok := os.LookupEnv("TOKEN"); ok {
		configs.Token = val
	}
	c, _ := w3s.NewClient(w3s.WithToken(configs.Token))
	if len(location) == 0 {
		err := ytdl.Download(ytdlPath, ytdlOptions)
		if err != nil {
			return err, errorYTDL
		}
		if configs.Token == "" {
			return err, errorIPFS
		}
		f, err := os.Open(filename)
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return err, errorYTDL
			} else {
				// Really kinda a hacky workaround to "WARNING: Requested formats are incompatible for merge and will be merged into mkv."
				spl := strings.Split(filename, ".")
				spl = append(spl[:len(spl)-1], "mkv")
				filename = strings.Join(spl, ".")
				f, _ = os.Open(filename)
			}
		}

		fmt.Printf("[ipfs] attempting to put file '%s'\n", filename)
		cid, err = c.Put(context.Background(), f)
		if err != nil {
			return err, errorIPFS
		}
		location = cid.String()
		err = table.Store(db, id, format, location)
		if err != nil {
			return err, errorTable
		}
	} else {
		fmt.Fprintf(stdout, "[w3s] Getting %s\n", location)
		res, _ := c.Get(context.Background(), cid)

		// res is a http.Response with an extra method for reading IPFS UnixFS files!
		f, fsys, _ := res.Files()

		// Download directory contents
		if d, ok := f.(fs.ReadDirFile); ok {
			ents, _ := d.ReadDir(0)
			for _, ent := range ents {
				file, err := fsys.Open("/" + ent.Name())
				if err != nil {
					return err, errorIPFS
				}
				data, err := ioutil.ReadAll(file)
				if err != nil {
					return err, errorIPFS
				}
				if err := ioutil.WriteFile(filename, data, 0755); err != nil {
					return err, errorIPFS
				}
			}
		}
	}

	fmt.Fprintf(stdout, "File is available locally at %s\n", filename)
	fmt.Fprintf(stdout, "File is also available at https://%s.ipfs.dweb.link\n", location)
	return nil, 0
}
