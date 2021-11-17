package ytdl

import (
	"log"
	"os/exec"
	"strings"
	"unicode/utf8"
)

// GetIdentifiers gets the identifying characteristics of the file to be downloaded
func GetIdentifiers(logger *log.Logger, binary string, args []string) (error, string, string) {
	logger.SetPrefix("[ytdl] ")
	logger.Println("Getting identifiers...")
	err, id := readCommand(binary, append(args, "--get-id"))
	if err != nil {
		return err, "", ""
	}
	err, format := readCommand(binary, append(args, "--get-format"))
	if err != nil {
		return err, "", ""
	}
	logger.Printf("Grabbed identifiers '%s' and '%s'\n", id, format)

	return nil, id, format
}

func GetFilename(logger *log.Logger, binary string, args []string) (error, string) {
	logger.SetPrefix("[ytdl] ")
	logger.Println("Getting filename")
	err, filename := readCommand(binary, append(args, "--get-filename"))
	if err != nil {
		return err, ""
	}
	spl := strings.Split(filename, ".")
	filename = strings.Join(spl[:len(spl)-1], ".")
	logger.Printf("Got filename %s\n", filename)
	return nil, filename
}

func Download(logger *log.Logger, binary string, args []string) error {
	logger.SetPrefix("[ytdl] ")
	err, _ := readCommand(binary, args)
	if err != nil {
		return err
	}
	logger.Println("File saved!")
	return nil
}

func readCommand(binary string, args []string) (error, string) {
	out, err := exec.Command(binary, append(args, "-o", "%(id)s -- %(format_id)s.%(ext)s")...).Output()
	if err != nil {
		return err, ""
	}

	strOut := string(out)
	r, size := utf8.DecodeLastRuneInString(strOut)
	if r == utf8.RuneError && (size == 0 || size == 1) {
		size = 0
	}
	return nil, string(out[:len(strOut)-size])
}
