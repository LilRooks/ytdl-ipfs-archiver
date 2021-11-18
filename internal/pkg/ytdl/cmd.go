package ytdl

import (
	"log"
	"os/exec"
	"strings"
	"unicode/utf8"
)

// GetIdentifiers gets the identifying characteristics of the file to be downloaded
func GetIdentifiers(logger *log.Logger, binary string, args []string) (string, string, error) {
	logger.SetPrefix("[ytdl] ")
	logger.Println("Getting identifiers...")
	id, err := readCommand(binary, append(args, "--get-id"))
	if err != nil {
		return "", "", err
	}
	format, err := readCommand(binary, append(args, "--get-format"))
	if err != nil {
		return "", "", err
	}
	logger.Printf("Grabbed identifiers '%s' and '%s'\n", id, format)

	return id, format, nil
}

// GetFilename determines the filename to write to
func GetFilename(logger *log.Logger, binary string, args []string) (string, error) {
	logger.SetPrefix("[ytdl] ")
	logger.Println("Getting filename")
	filename, err := readCommand(binary, append(args, "--get-filename"))
	if err != nil {
		return "", err
	}
	spl := strings.Split(filename, ".")
	filename = strings.Join(spl[:len(spl)-1], ".")
	logger.Printf("Got filename %s\n", filename)
	return filename, nil
}

// Download executes the assumed download file
func Download(logger *log.Logger, binary string, args []string) error {
	logger.SetPrefix("[ytdl] ")
	_, err := readCommand(binary, args)
	if err != nil {
		return err
	}
	logger.Println("File saved!")
	return nil
}

func readCommand(binary string, args []string) (string, error) {
	out, err := exec.Command(binary, append(args, "-o", "%(id)s -- %(format_id)s.%(ext)s")...).Output()
	if err != nil {
		return "", err
	}

	strOut := string(out)
	r, size := utf8.DecodeLastRuneInString(strOut)
	if r == utf8.RuneError && (size == 0 || size == 1) {
		size = 0
	}
	return string(out[:len(strOut)-size]), nil
}
