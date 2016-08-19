package main

import (
	"fmt"
	"os"

	video2hevc "github.com/martinlindhe/video2hevc/lib"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	file = kingpin.Arg("file", "Input file").Required().String()
)

func main() {

	// support -h for --help
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Parse()

	err := video2hevc.VideoToHevc(*file)
	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}
}
