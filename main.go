package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

var PATHS []string
var REQUIRED_FILES = make(map[string]*SourceFile)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "src",
				Usage:    "The main, top-level source file",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "dest",
				Usage:    "The output file",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "path",
				Usage: "Additional paths to search for files, split by semicolon (;)",
			},
		},
		Action: func(cCtx *cli.Context) error {
			src := cCtx.String("src")
			dest := cCtx.String("dest")
			path := cCtx.String("path")

			// "C:/Users/micro/AppData/Roaming/com.nesbox.tic/TIC-80/rift-nova-2023/demo.lua", "./out.lua"
			doBundling(src, dest, path)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func doBundling(srcMainFilepath string, destFilepath string, pathsString string) {
	if pathsString != "" {
		paths := strings.Split(pathsString, ";")
		for _, path := range paths {
			strings.Trim(path, " ")
			if path != "" {
				addPath(path)
			}
		}
	}

	file, err := openFile(srcMainFilepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	mainSource, err := readFileLines(file)
	if err != nil {
		log.Fatal(err)
	}

	mainFile := NewSourceFile("main", mainSource)

	// Read In
	mainFile.Strip()

	// Write all
	OutFile, err := os.Create(destFilepath)
	if err != nil {
		log.Fatal(err)
	}
	defer OutFile.Close()

	writeLine(OutFile, "-- Assembled by the RiFT bundler\n\n")

	// Ensure all files are stripped.
	// We're adding to this map as we iterate over it - so this is over the top, but working
	for {
		doAgain := false
		for _, file := range REQUIRED_FILES {
			if !file.isStripped {
				file.Strip()
				// Might have added some new sources
				doAgain = true
			}
		}

		if !doAgain {
			break
		}
	}

	// Write out the files
	for funcName, rf := range REQUIRED_FILES {
		fmt.Printf("Bundling file: %s\n", funcName)
		writeLine(OutFile, fmt.Sprintf("%s=function()\n", funcName))
		rf.Write(OutFile)
		writeLine(OutFile, "\nend\n\n")
	}

	mainFile.Write(OutFile)

	fmt.Println("-> DONE. Hope it works!")
}

func readFile(filename string) string {
	body, err := os.ReadFile("file.txt")
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}
	return string(body)
}

func openFile(path string) (*os.File, error) {
	return os.Open(path)
}

func readFileLines(file *os.File) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func writeLine(file *os.File, line string) {
	_, err := file.Write([]byte(line))
	if err != nil {
		log.Fatal(err)
	}
}
