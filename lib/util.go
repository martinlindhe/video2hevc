package video2hevc

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

type VideoToHevcSettings struct {
	ForceAAC    bool
	ForceAC3    bool
	ForceNVIDIA bool
	Verbose     bool
	Force720    bool
	Threads     int
}

// VideoToHevc encodes video `file` using command line `ffmpeg` tool
func VideoToHevc(file string, settings VideoToHevcSettings) error {
	if !exists(file) {
		return fmt.Errorf("%s not found", file)
	}

	outName := findFreeOutFileName(file)

	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		return err
	}
	audioLib := "copy"
	if settings.ForceAAC {
		audioLib = "aac"
	}
	if settings.ForceAC3 {
		audioLib = "ac3"
	}
	videoLib := "libx265"
	if settings.ForceNVIDIA {
		videoLib = "hevc_nvenc"
	}

	parameters := []string{
		"-i", file,
	}

	if settings.Threads != 0 {
		parameters = append(parameters, []string{"-threads", fmt.Sprintf("%d", settings.Threads)}...)
	}
	parameters = append(parameters, []string{
		"-c:v", videoLib,
		"-c:a", audioLib,
	}...)

	if settings.Force720 {
		// 1280 x 720
		parameters = append(parameters, []string{"-vf", "scale=-1:720"}...)
	}

	if settings.Verbose {
		parameters = append(parameters, []string{"-loglevel", "verbose"}...)
	}

	parameters = append(parameters, outName)

	if settings.Verbose {
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
