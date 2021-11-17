package ipfs

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/web3-storage/go-w3s-client"
)

func Store(c w3s.Client, filename string) (error, string) {
	f, err := os.Open(filename)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err, ""
		} else {
			// Really kinda a hacky workaround to "WARNING: Requested formats are incompatible for merge and will be merged into mkv."
			spl := strings.Split(filename, ".")
			spl = append(spl[:len(spl)-1], "mkv")
			filename = strings.Join(spl, ".")
			f, _ = os.Open(filename)
		}
	}
	fmt.Printf("[ipfs] attempting to put file '%s'\n", filename)
	cid, err := c.Put(context.Background(), f)
	if err != nil {
		return err, ""
	}
	return nil, cid.String()
}
