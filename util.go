package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

func baseNameWithoutExt(filename string) string {

	s := filepath.Base(filename)
	n := strings.LastIndexByte(s, '.')
	if n >= 0 {
		return s[:n]
	}
	return s
}

// exists reports whether the named file or directory exists.
func exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func findFreeOutFileName(file string) string {

	cnt := 0
	res := ""
	for {
		res = path.Join(filepath.Dir(file), baseNameWithoutExt(file))
		if cnt > 0 {
			res += "-" + fmt.Sprintf("%02d", cnt)
		}
		res += ".mp4"
		if !exists(res) {
			break
		}
		cnt++
	}
	return res
}

// interactive commands (ssh, vim)
func runInteractiveCommand(name string, arg ...string) error {

	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func runCommandReturnStdout(step string) (string, error) {

	parts := strings.Split(step, " ")
	cmd := exec.Command(parts[0], parts[1:]...)
	res := ""

	stdOutReader, err := cmd.StdoutPipe()
	if err != nil {
		return res, err
	}
	stdOutScanner := bufio.NewScanner(stdOutReader)
	go func() {
		for stdOutScanner.Scan() {
			res += string(stdOutScanner.Bytes()) + "\n"
		}
	}()

	if err := cmd.Start(); err != nil {
		return res, err
	}
	if err := cmd.Wait(); err != nil {
		return res, err
	}
	return res, nil
}
