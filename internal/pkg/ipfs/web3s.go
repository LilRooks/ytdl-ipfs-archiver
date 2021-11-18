package ipfs

import (
	"context"
	"fmt"
	"os"

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

// Fetch grabs the file referred by the cid
func Fetch(c w3s.Client, cidStr string) error {
	// TODO fetching functionality
	return nil
}
