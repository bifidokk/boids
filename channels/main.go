package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	windRegex        = regexp.MustCompile(`\d* METAR.*EGLL \d*Z [A-Z ]*(\d{5}KT|VRB\d{2}KT).*=`)
	tafValidation    = regexp.MustCompile(`.*TAF.*`)
	comment          = regexp.MustCompile(`\w*#.*`)
	metarClose       = regexp.MustCompile(`.*=`)
	variableWind     = regexp.MustCompile(`.*VRB\d{2}KT`)
	validWind        = regexp.MustCompile(`\d{5}KT`)
	windDirOnly      = regexp.MustCompile(`(\d{3})\d{2}KT`)
	windDistribution [8]int
)

func main() {
	textChannel := make(chan string)
	metarChannel := make(chan []string)
	windsChannel := make(chan []string)
	resultsChannel := make(chan [8]int)

	go parseToArray(textChannel, metarChannel)
	go extractWindDirection(metarChannel, windsChannel)
	go mineWindDistribution(windsChannel, resultsChannel)

	abspath, _ := filepath.Abs("./metarfiles/")
	files, _ := os.ReadDir(abspath)
	start := time.Now()

	for _, file := range files {
		data, err := os.ReadFile(filepath.Join(abspath, file.Name()))
		if err != nil {
			panic(err)
		}

		text := string(data)
		textChannel <- text
	}

	close(textChannel)

	results := <-resultsChannel

	elapsed := time.Since(start)
	fmt.Printf("Took %s\n", elapsed)
	fmt.Printf("%v\n", results)
}

func parseToArray(textChannel chan string, metarChannel chan []string) {
	for text := range textChannel {
		lines := strings.Split(text, "\n")
		metarSlice := make([]string, 0, len(lines))
		metarStr := ""

		for _, line := range lines {
			if tafValidation.MatchString(line) {
				break
			}

			if !comment.MatchString(line) {
				metarStr += strings.Trim(line, " ")
			}

			if metarClose.MatchString(line) {
				metarSlice = append(metarSlice, metarStr)
				metarStr = ""
			}
		}

		metarChannel <- metarSlice
	}

	close(metarChannel)
}

func extractWindDirection(metarChannel chan []string, windsChannel chan []string) {
	for metars := range metarChannel {
		winds := make([]string, 0, len(metars))

		for _, metar := range metars {
			if windRegex.MatchString(metar) {
				winds = append(winds, windRegex.FindAllStringSubmatch(metar, -1)[0][1])
			}
		}

		windsChannel <- winds
	}

	close(windsChannel)
}

func mineWindDistribution(windsChannel chan []string, resultsChannel chan [8]int) {
	for winds := range windsChannel {
		for _, wind := range winds {
			if variableWind.MatchString(wind) {
				for i := 0; i < 8; i++ {
					windDistribution[i]++
				}
			} else if validWind.MatchString(wind) {
				windStr := windDirOnly.FindAllStringSubmatch(wind, -1)[0][1]

				if d, err := strconv.ParseFloat(windStr, 64); err == nil {
					dirIndex := int(math.Round(d/45.0)) % 8
					windDistribution[dirIndex]++
				}
			}
		}
	}

	resultsChannel <- windDistribution
	close(resultsChannel)
}
