package main

import (
	"fmt"
	"os"

	video2hevc "github.com/martinlindhe/video2hevc/lib"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	file    = kingpin.Arg("file", "Input file").Required().String()
	aac     = kingpin.Flag("aac", "Force AAC audio (better quality then ac3)").Bool()
	ac3     = kingpin.Flag("ac3", "Force AC3 audio").Bool()
	nvidia  = kingpin.Flag("nvidia", "Force NVIDIA acceleration").Bool()
	verbose = kingpin.Flag("verbose", "Be verbose").Short('v').Bool()
	v720    = kingpin.Flag("v720", "Convert video to 720p").Bool()
	threads = kingpin.Flag("threads", "Number of threads").Default("0").Int()
)

func main() {

	// support -h for --help
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Parse()

	settings := video2hevc.VideoToHevcSettings{
		ForceAAC:    *aac,
		ForceAC3:    *ac3,
		ForceNVIDIA: *nvidia,
		Verbose:     *verbose,
		Force720:    *v720,
		Threads:     *threads,
	}

	err := video2hevc.VideoToHevc(*file, settings)
	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}
}
