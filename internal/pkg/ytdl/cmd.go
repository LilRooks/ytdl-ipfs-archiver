package ytdl

import (
	"fmt"
	"os/exec"
	"unicode/utf8"
)

// GetIdentifiers gets the identifying characteristics of the file to be downloaded
func GetIdentifiers(binary string, args []string) (error, string, string) {
	err, id := ReadCommand(binary, append(args, "--get-id"))
	if err != nil {
		return err, "", ""
	}
	err, format := ReadCommand(binary, append(args, "--get-format"))
	if err != nil {
		return err, "", ""
	}
	fmt.Printf("[ytdl] Grabbed identifiers '%s' and '%s'\n", id, format)
	return nil, id, format
}

func GetFilename(binary string, args []string) (error, string) {
	err, filename := ReadCommand(binary, append(args, "--get-filename"))
	if err != nil {
		return err, ""
	}
	fmt.Println("[ytdl] Got filename", filename)
	return nil, filename
}

func Download(binary string, args []string) error {
	err, _ := ReadCommand(binary, args)
	if err != nil {
		return err
	}
	fmt.Println("[ytdl] file saved!")
	return nil
}

func ReadCommand(binary string, args []string) (error, string) {
	out, err := exec.Command(binary, args...).Output()
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
