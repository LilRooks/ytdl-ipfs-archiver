package ytdl

import "os/exec"

// ParsePath parses the provided path into an os/exec usable binary path
func ParsePath(pathArg string) (error, string) {
	if pathArg == "" {
		pathArg = "youtube-dl"
	}
	ytdlExecPath, err := exec.LookPath(pathArg)
	if err != nil {
		return err, ""
	} else {
		return nil, ytdlExecPath
	}
	return nil, pathArg
}
