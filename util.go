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

func videoToHevc(file string) error {

	if !exists(file) {
		return fmt.Errorf("%s not found", file)
	}

	outName := findFreeOutFileName(file)

	// find ffmpeg exec path, works on macos, debian
	ffmpegPath, err := runCommandReturnStdout("/usr/bin/which ffmpeg")
	if err != nil {
		return fmt.Errorf("couldn't call which: %s", err)
	}
	ffmpegPath = strings.TrimSpace(ffmpegPath)
	if ffmpegPath == "" {
		return fmt.Errorf("could not find ffmpeg binary")
	}

	arg := []string{"-i", file, "-c:v", "libx265", "-c:a", "libfdk_aac", outName}
	fmt.Println("cmd:", ffmpegPath, strings.Join(arg, " "))

	err = runInteractiveCommand(ffmpegPath, arg...)
	if err != nil {
		return fmt.Errorf("exec error: %s", err)
	}
	return nil
}

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
