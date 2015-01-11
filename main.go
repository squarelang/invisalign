package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Counts the instances of the " " character that prefix the given string.

func countLeadingSpaces(s string) int {
	if len(s) == 0 || !strings.HasPrefix(s, " ") {
		return 0
	}
	return 1 + countLeadingSpaces(s[1:len(s)])
}

// Adds curly braces based on symantic whitespacing.

func addBraces(lines []string) (linesOut []string) {
	// This orthodontist puts the braces ON
	lines = append(lines, "")
	prevLine := struct{ spaces, lineNum int }{}
	spaceCount := []int{}

	for n, line := range lines {
		spaces := countLeadingSpaces(line)
		currLine := struct{ spaces, lineNum int }{spaces, n}
		spaceCount = append(spaceCount, currLine.spaces)
		if line == "" && n < len(lines)-1 {
			// Pass through on blank lines
			linesOut = append(linesOut, line)
			continue
		}
		if prevLine.spaces < currLine.spaces {
			// This line is indented, put a { on the previous nontrivial line
			linesOut[prevLine.lineNum] += " {"
		}
		if prevLine.spaces > currLine.spaces || n == len(lines)-1 {
			// This line is outdented, walk upwards to see how many closing }'s to place
			tempSpaces := prevLine.spaces
			linesOut[prevLine.lineNum] += " }"
			for i := n - 1; lines[i] == "" || spaceCount[i] > spaceCount[n]; i-- {
				if lines[i] != "" && tempSpaces > spaceCount[i] {
					linesOut[prevLine.lineNum] += " }"
					tempSpaces = spaceCount[i]
				}
			}
		}
		linesOut = append(linesOut, line)
		prevLine = currLine
	}
	return
}

// Reads a file to a slice of strings.

func readFileToSlice(filename string) (lines []string, err error) {
	file, err := os.Open(filename)
	defer file.Close()
	if err == nil {
		fileScanner := bufio.NewScanner(file)
		for fileScanner.Scan() {
			lines = append(lines, fileScanner.Text())
		}
	}
	return
}

// Invisalign is the full compile function. Accepts a filepath and an optional parameter '-w' which,
// if given, says to overwrite the given file with the compiled go. Otherwise, prints the go code to
// stdout.

func invisalign(args []string) {
	filenameIn := args[len(args)-1]

	// Parse source file
	linesIn, err := readFileToSlice(filenameIn)
	if err != nil {
		log.Fatalf("Couldn't open input file for reading: " + err.Error())
	}

	// Visit the orthodontist
	linesOut := addBraces(linesIn)

	// Write to temp file
	output, err := ioutil.TempFile("/tmp", filepath.Base(filenameIn))
	if err != nil {
		log.Fatalf("Couldn't open temporary file for writing: " + err.Error())
	}
	for _, line := range linesOut {
		output.WriteString(line + "\n")
	}
	output.Close()

	// Gofmt
	args[len(args)-1] = output.Name()
	cmd := exec.Command("gofmt", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	// Overwrite original file
	if args[0] == "-w" {
		os.Rename(output.Name(), filenameIn)
	}
}

func main() {
	args := os.Args[1:len(os.Args)]
	invisalign(args)
}
