package web

import (
	"bufio"
	"bytes"
	"os"
	"strings"
)

const ENV_FILE_NAME = ".env"

func OverWriteEnvFile(target map[string]string) error {
	f, err := os.Open(ENV_FILE_NAME)
	if err != nil {
		return err
	}
	defer f.Close()
	var bs []byte
	buf := bytes.NewBuffer(bs)
	scanner := bufio.NewScanner(f)
	flag := false
	for scanner.Scan() {
		for t := range target {
			if strings.Contains(scanner.Text(), t) {
				_, err := buf.WriteString(t + "=" + target[t] + "\n")
				if err != nil {
					return err
				}
				flag = true
				break
			}
		}
		if flag == true {
			flag = false
			continue
		}
		_, err = buf.WriteString(string(scanner.Bytes()) + "\n")
		if err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	err = os.WriteFile(ENV_FILE_NAME, buf.Bytes(), 0666)
	if err != nil {
		return err
	}
	return nil
}
