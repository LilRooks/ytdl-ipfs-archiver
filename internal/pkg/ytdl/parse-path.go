package ytdl

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
)

//go:embed ytdl
var ytdl []byte

// ParsePath parses the provided path into an os/exec usable binary path
func ParsePath(pathArg string) (error, string) {
	if len(pathArg) == 0 {
		pathArg = "youtube-dl"
	} else if pathArg == "embedded" {
		var err error
		err, pathArg = embeddedYtdl()
		if err != nil {
			return err, ""
		}
	}
	ytdlExecPath, err := exec.LookPath(pathArg)
	if err != nil {
		return err, ""
	} else {
		fmt.Printf("[ytdl] found binary at '%s'\n", ytdlExecPath)
		return nil, ytdlExecPath
	}
	return nil, pathArg
}

func embeddedYtdl() (error, string) {
	// Create and write to file
	tmpFile, err := os.CreateTemp(os.TempDir(), "ytdl-*")
	if err != nil {
		return err, ""
	}
	if _, err := tmpFile.Write(ytdl); err != nil {
		return err, ""
	}
	ytdlExecPath := tmpFile.Name()
	// Executable
	if err := tmpFile.Chmod(0755); err != nil {
		return err, ""
	}
	if err := tmpFile.Close(); err != nil {
		return err, ""
	}
	return nil, ytdlExecPath
}
