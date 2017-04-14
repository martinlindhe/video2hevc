package main

import (
	"fmt"
	"os"

	video2hevc "github.com/martinlindhe/video2hevc/lib"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	file   = kingpin.Arg("file", "Input file").Required().String()
	aac    = kingpin.Flag("aac", "Force AAC audio").Bool()
	nvidia = kingpin.Flag("nvidia", "Force NVIDIA acceleration").Bool()
)

func main() {

	// support -h for --help
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Parse()

	err := video2hevc.VideoToHevc(*file, *aac, *nvidia)
	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}
}
