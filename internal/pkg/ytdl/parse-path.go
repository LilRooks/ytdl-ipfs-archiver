package ytdl

import (
	"fmt"
	"os/exec"
)

// ParsePath parses the provided path into an os/exec usable binary path
func ParsePath(pathArg string) (string, error) {
	if len(pathArg) == 0 {
		pathArg = "youtube-dl"
	}
	ytdlExecPath, err := exec.LookPath(pathArg)
	if err != nil {
		return "", err
	}
	fmt.Printf("[ytdl] found binary at '%s'\n", ytdlExecPath)
	return ytdlExecPath, nil
}
