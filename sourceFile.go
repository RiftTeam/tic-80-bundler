package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

var REGEX_PACKAGE_PATH = regexp.MustCompile(`package.path\w*=\w*package.path\w*..\w*\"([;\w:\\.\-\?]+)\"`)
var REGEX_REQUIRE = regexp.MustCompile(`(.*)require\([\"']([\w\._\\\-\/]+)[\"']\)(\([\w,\s]+\))?(.*)`)

type SourceFile struct {
	name         string
	code         []string
	isStripped   bool
	codeStripped []string
}

func NewSourceFile(name string, code []string) *SourceFile {
	return &SourceFile{name: name, code: code, isStripped: false, codeStripped: make([]string, 0)}
}

func (f *SourceFile) IsStripped() bool {
	return f.isStripped
}

func (f *SourceFile) Strip() {
	for _, line := range f.code {
		includeLine := true
		matchesPath := REGEX_PACKAGE_PATH.FindStringSubmatch(line)
		if len(matchesPath) > 0 {
			decodeAndAddPath(matchesPath[1])
			includeLine = false
		}

		line = ReplaceAllStringSubmatchFunc(REGEX_REQUIRE, line, func(groups []string) string {
			filename := groups[2]
			funcName := addSourceFile(filename)

			return fmt.Sprintf("%s%s()%s%s", groups[1], funcName, groups[3], groups[4])
		})

		if includeLine {
			f.codeStripped = append(f.codeStripped, line)
		}
	}

	f.isStripped = true
}

func (f *SourceFile) Write(outFile *os.File) {
	fmt.Printf("Lines Out: %d (%s)\n", len(f.codeStripped), f.name)
	for _, line := range f.codeStripped {
		outFile.WriteString(line + "\n")
	}
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

func getFuncName(filename string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	return "rift_" + reg.ReplaceAllString(filename, "")
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
