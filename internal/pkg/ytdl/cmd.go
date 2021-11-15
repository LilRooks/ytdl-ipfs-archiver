package ytdl

import (
	"log"
	"os/exec"
	"unicode/utf8"
)

func GetIdentifiers(binary string, args []string) (error, string, string) {
	err, id := ReadCommand(binary, append(args, "--get-id"))
	if err != nil {
		return err, "", ""
	}
	err, format := ReadCommand(binary, append(args, "--get-format"))
	if err != nil {
		return err, "", ""
	}

	return nil, id, format
}

func ReadCommand(binary string, args []string) (error, string) {
	out, err := exec.Command(binary, args...).Output()
	if err != nil {
		log.Fatal(err)
	}

	strOut := string(out)
	r, size := utf8.DecodeLastRuneInString(strOut)
	if r == utf8.RuneError && (size == 0 || size == 1) {
		size = 0
	}
	return nil, string(out[:len(strOut)-size])
}
