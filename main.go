package main

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	file = kingpin.Arg("file", "Input file").Required().String()
)

func main() {

	// support -h for --help
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Parse()

	if !exists(*file) {
		fmt.Println("error: input file not found")
		os.Exit(1)
	}

	outName := findFreeOutFileName(*file)

	// find ffmpeg exec path, works on macos, debian
	ffmpegPath, err := runCommandReturnStdout("/usr/bin/which ffmpeg")
	if err != nil {
		fmt.Println("which err:", err)
		os.Exit(1)
	}
	ffmpegPath = strings.TrimSpace(ffmpegPath)

	arg := []string{"-i", *file, "-c:v", "libx265", "-c:a", "libfdk_aac", outName}
	fmt.Println("cmd:", ffmpegPath, strings.Join(arg, " "))

	err = runInteractiveCommand(ffmpegPath, arg...)
	if err != nil {
		fmt.Println("exec error:", err)
	}
}
