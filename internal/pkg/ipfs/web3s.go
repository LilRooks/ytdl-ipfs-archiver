package ipfs

import (
	"context"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"

	"github.com/ipfs/go-cid"
	"github.com/web3-storage/go-w3s-client"
)

// Store stores a file to web3.storage
func Store(c w3s.Client, filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	fmt.Printf("[ipfs] attempting to put file '%s'\n", filename)
	cid, err := c.Put(context.Background(), f)
	if err != nil {
		return "", err
	}
	return cid.String(), nil
}

// Fetch grabs the file referred by the cid, returns local path
func Fetch(c w3s.Client, cidStr string) (string, error) {
	fmt.Printf("[ipfs] attempting to pull file '%s'\n", cidStr)
	var filename string
	cid, err := cid.Decode(cidStr)
	if err != nil {
		return "", err
	}
	res, err := c.Get(context.Background(), cid)
	if err != nil {
		return "", err
	}

	// Download directory contents
	f, fsys, err := res.Files()
	if err != nil {
		return "", err
	}
	if d, ok := f.(fs.ReadDirFile); ok {
		ents, err := d.ReadDir(0)
		if err != nil {
			return "", err
		}
		for _, ent := range ents {
			filename = ent.Name()
			file, err := fsys.Open("/" + ent.Name())
			if err != nil {
				return "", err
			}

			data, err := ioutil.ReadAll(file)
			if err != nil {
				return "", err
			}

			err = ioutil.WriteFile(ent.Name(), data, 0755)
			if err != nil {
				return "", err
			}
		}
	} // TODO fetching functionality
	return filename, nil
}
