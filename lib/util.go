package video2hevc

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

// VideoToHevc encodes video `file` using command line `ffmpeg` tool
func VideoToHevc(file string, verbose bool, forceAAC bool, forceAC3 bool, forceNVIDIA bool, v720 bool) error {
	if !exists(file) {
		return fmt.Errorf("%s not found", file)
	}

	outName := findFreeOutFileName(file)

	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		return err
	}
	audioLib := "copy"
	if forceAAC {
		audioLib = "aac"
	}
	if forceAC3 {
		audioLib = "ac3"
	}
	videoLib := "libx265"
	if forceNVIDIA {
		videoLib = "hevc_nvenc"
	}

	parameters := []string{
		"-i", file,
		"-c:v", videoLib,
		"-c:a", audioLib,
	}
	if v720 {
		// 1280 x 720
		parameters = append(parameters, []string{"-vf", "scale=-1:720"}...)
	}
	if verbose {
		parameters = append(parameters, []string{"-loglevel", "verbose"}...)
	}

	parameters = append(parameters, outName)

	if verbose {
		fmt.Println("Executing", ffmpegPath, strings.Join(parameters, " "))
	}
	err = runInteractiveCommand(ffmpegPath, parameters...)
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
	ext := ".mp4"
	for {
		res = path.Join(filepath.Dir(file), baseNameWithoutExt(file))
		if cnt > 0 {
			res += "-" + fmt.Sprintf("%02d", cnt)
		}
		res += ext
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

	return cmd.Run()
}
