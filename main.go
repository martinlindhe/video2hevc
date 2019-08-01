package main

import (
	"fmt"
	"os"

	video2hevc "github.com/martinlindhe/video2hevc/lib"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	file    = kingpin.Arg("file", "Input file").Required().String()
	aac     = kingpin.Flag("aac", "Force AAC audio (default copies audio)").Bool()
	ac3     = kingpin.Flag("ac3", "Force AC3 audio (default copies audio)").Bool()
	nvidia  = kingpin.Flag("nvidia", "Force NVIDIA acceleration").Bool()
	verbose = kingpin.Flag("verbose", "Be verbose").Short('v').Bool()
	v720    = kingpin.Flag("v720", "Convert video yo 720p").Bool()
)

func main() {

	// support -h for --help
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Parse()

	err := video2hevc.VideoToHevc(*file, *verbose, *aac, *ac3, *nvidia, *v720)
	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}
}
