package main

import (
	"fmt"
	"os"
)

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
	fmt.Printf("Lines In: %d (%s)\n", len(f.code), f.name)
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
