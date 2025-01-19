package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Dir struct {
	Dir     string   `json:"dir"`
	Files   []string `json:"files"`
	Folders []Dir    `json:"folders"`
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	line, _ := reader.ReadString('\n')
	t, _ := strconv.Atoi(strings.TrimSpace(line))

	for i := 0; i < t; i++ {
		line, _ = reader.ReadString('\n')
		n, _ := strconv.Atoi(strings.TrimSpace(line))

		var jsonBlob []byte

		for j := 0; j < n; j++ {
			jsonLine, _ := reader.ReadString('\n')
			jsonBlob = append(jsonBlob, []byte(strings.ReplaceAll(strings.TrimSpace(jsonLine), "\n", ""))...)
		}

		var dir Dir

		err := json.Unmarshal(jsonBlob, &dir)

		if err != nil {
			panic(err)
		}

		infectedCount := countInfectedFiles(dir, false)
		fmt.Fprintln(writer, infectedCount)
	}
}

func countInfectedFiles(dir Dir, parentInfected bool) int {
	infectedCount := 0
	hasVirus := parentInfected

	for _, file := range dir.Files {
		if isInfected(file) {
			hasVirus = true
		}
	}

	// if current folder is infected, then all nested folders and files are infected too
	if hasVirus {
		infectedCount += len(dir.Files)
	}

	for _, folder := range dir.Folders {
		infectedCount += countInfectedFiles(folder, hasVirus)
	}

	return infectedCount
}

func isInfected(file string) bool {
	return strings.HasSuffix(file, ".hack")
}

/*
Files and folders are considered infected if they are in the same directory as the virus. If a folder is infected, it means that all files and folders inside are also infected. A virus file is considered an infected file. All virus files have the extension .hack For example: cat.png.hack , a.cpp.hack are considered a virus, config.yaml, hack.cpp are not considered a virus.

Input
2
23
{
"dir": "root",
"files": [".zshrc"],
"folders": [
{
"dir": "desktop",
"files": ["config.yaml"]
},
{
"dir": "downloads",
"files": ["cat.png.hack"],
"folders": [
{
"dir": "kta",
"files": [
"kta.exe",
"kta.hack"
]
}
]
}
]
}
8
{ "dir"
:
"awesomeproject"
,
"files" : [ "go.mod",
"go.sum"
],
"folders":[{"dir":"cmd"}]}


Output
3
0
*/
