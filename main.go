package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

var PATHS []string
var MAIN_FILE *SourceFile
var REQUIRED_FILES = make(map[string]*SourceFile)
var OUTFILE *os.File

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
		},
		Action: func(cCtx *cli.Context) error {
			src := cCtx.String("src")
			dest := cCtx.String("dest")
			// "C:/Users/micro/AppData/Roaming/com.nesbox.tic/TIC-80/rift-nova-2023/demo.lua", "./out.lua"
			doBundling(src, dest)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func doBundling(srcMainFilepath string, destFilepath string) {
	file, err := openFile(srcMainFilepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	mainSource, err := readFileLines(file)
	MAIN_FILE = NewSourceFile("main", mainSource)

	// Read In
	MAIN_FILE.Strip()

	// Write all
	OUTFILE, err := os.Create(destFilepath)
	if err != nil {
		log.Fatal(err)
	}
	defer OUTFILE.Close()

	writeLine(OUTFILE, "-- Assembled by the RiFT bundler\n\n")

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
		fmt.Printf("Writing file: %s\n", funcName)
		writeLine(OUTFILE, fmt.Sprintf("%s=function()\n", funcName))
		rf.Write(OUTFILE)
		writeLine(OUTFILE, "\nend\n\n")
	}

	MAIN_FILE.Write(OUTFILE)

	fmt.Println("done")
}

func readFile(filename string) string {
	body, err := ioutil.ReadFile("file.txt")
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
