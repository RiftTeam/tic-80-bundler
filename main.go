package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

const IN_FILENAME = "C:/Users/micro/AppData/Roaming/com.nesbox.tic/TIC-80/rift-nova-2023/demo.lua"
const OUT_FILENAME = "./out.lua"

var PATHS []string
var REGEX_PACKAGE_PATH = regexp.MustCompile(`package.path\w*=\w*package.path\w*..\w*\"([;\w:\\.\-\?]+)\"`)
var REGEX_REQUIRE = regexp.MustCompile(`(.*)require\([\"']([\w\._\\\-\/]+)[\"']\)(\([\w,\s]+\))?(.*)`)
var MAIN_FILE *SourceFile
var REQUIRED_FILES = make(map[string]*SourceFile)
var OUTFILE *os.File

func main() {
	file, err := openFile(IN_FILENAME)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	mainSource, err := readFileLines(file)
	MAIN_FILE = NewSourceFile("main", mainSource)

	// Read In
	MAIN_FILE.Strip()

	// Write all
	OUTFILE, err := os.Create(OUT_FILENAME)
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

func decodeAndAddPath(path string) {
	path = filepath.FromSlash(path[1 : len(path)-5])
	if stat, err := os.Stat(path); err == nil && stat.IsDir() {
		fmt.Printf("Path (added): %s\n", path)
		PATHS = append(PATHS, path)
	} else {
		fmt.Printf("Path (ignored): %s\n", path)
	}
}

func addSourceFile(filename string) string {
	funcName := getFuncName(filename)
	if _, ok := REQUIRED_FILES[funcName]; ok {
		return funcName // We already have this one
	}

	file, err := findAndOpenFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	source, err := readFileLines(file)
	REQUIRED_FILES[funcName] = NewSourceFile(filename, source)

	return funcName
}

// Use the paths we've found so far
// ensure you call defer file.Close()
func findAndOpenFile(filename string) (*os.File, error) {
	for _, path := range PATHS {
		fullPath := filepath.Join(path, filename) + ".lua"
		fmt.Printf("Trying required file: %s\n", fullPath)

		file, err := os.Open(fullPath)
		if err == nil {
			return file, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("File not found (%s)", filename))
}

// PATHS
// package.path=package.path..";C:\\Users\\micro\\AppData\\Roaming\\com.nesbox.tic\\TIC-80\\rift\\?.lua"   -- jtruk
// require("./sys/sys")(R)

// require("./state-logo"),

func getFuncName(filename string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	return "rift_" + reg.ReplaceAllString(filename, "")
}

func readFile(filename string) string {
	body, err := ioutil.ReadFile("file.txt")
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}
	return string(body)
}

func areAllFilesStripped() bool {
	for _, file := range REQUIRED_FILES {
		if file.isStripped {
			return false
		}
	}
	return true
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

// From: https://gist.github.com/elliotchance/d419395aa776d632d897?permalink_comment_id=3713809#gistcomment-3713809
func ReplaceAllStringSubmatchFunc(re *regexp.Regexp, str string, repl func([]string) string) string {
	result := ""
	lastIndex := 0

	for _, v := range re.FindAllSubmatchIndex([]byte(str), -1) {
		groups := []string{}
		for i := 0; i < len(v); i += 2 {
			if v[i] == -1 || v[i+1] == -1 {
				groups = append(groups, "")
			} else {
				groups = append(groups, str[v[i]:v[i+1]])
			}
		}

		result += str[lastIndex:v[0]] + repl(groups)
		lastIndex = v[1]
	}

	return result + str[lastIndex:]
}

/*
First stab - might need to go back?

func getMagicCallFn() string {
	return "RIFT_CALLFN=function(fn,...)fn()(...)end\n\n"
}

func getWrappedCall(fn string, args string) string {
	return fmt.Sprintf("RIFT_CALLFN(function()\n%s\nend,%s)\n\n", fn, args)
}

/*
RIFT_CALLFN=function(fn,...)fn()(...)end

RIFT_CALLFN(function()
    -- Replace with your function
    return function(arg1,arg2)
        trace(arg1)
        trace(arg2)
    end
end,"the1st","and2nd")
*/
