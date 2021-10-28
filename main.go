package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sunshineplan/imgconv"
)

const photobucketRegex = `http:\/\/i[0-9]{2,}\.photobucket\.com\/albums\/[a-zA-Z]\d{2,}\/[^\/]{1,}\/.*?\.[a-z]{2,}`

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage ./photobucketConverter [inputFile]")
	}

	file, err := os.Open(os.Args[1])
	printAndQuitOnError(err)
	defer file.Close()

	outDir, err := os.MkdirTemp(".", "export_")
	printAndQuitOnError(err)

	scanner := bufio.NewScanner(file)
	compiledRegexp := regexp.MustCompile(photobucketRegex)

	var rawMatches []string
	for scanner.Scan() {
		rawMatches = append(rawMatches, compiledRegexp.FindAllString(scanner.Text(), -1)...)
		fmt.Println("%n", len(rawMatches))
	}
	matches := removeDuplicateStr(rawMatches)

	for i, match := range matches {
		fmt.Printf("Downloading file %d of %d: %s\n", i+1, len(matches), match)
		downloadAndWriteFile(match, outDir)
	}

	fmt.Printf("Data saved in %s\n\n", outDir)
}

func downloadAndWriteFile(url, outDir string) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "image/webp,*/*")
	req.Header.Add("Referer", url)

	resp, err := client.Do(req)
	printAndQuitOnError(err)

	if resp.StatusCode != 200 {
		fmt.Println("Download failed HTTP status " + resp.Status)
		return
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	printAndQuitOnError(err)

	urlSplits := strings.Split(url, "/")
	fileName := urlSplits[len(urlSplits)-1]
	fileExtension := strings.Split(fileName, ".")[len(strings.Split(fileName, "."))-1]

	if fileExtension == "webp" {
		err = os.WriteFile(filepath.Join(outDir, fileName), respBody, 0600)
		printAndQuitOnError(err)
	} else {
		format, err := imgconv.FormatFromExtension(fileExtension)
		printAndQuitOnError(err)

		buf := new(bytes.Buffer)
		img, err := imgconv.Decode(bytes.NewReader(respBody))
		err = imgconv.Write(buf, img, imgconv.FormatOption{Format: format})
		if err != nil {
			printAndQuitOnError(err)
		}
		printAndQuitOnError(err)
		err = os.WriteFile(filepath.Join(outDir, fileName), buf.Bytes(), 0600)
		printAndQuitOnError(err)
	}
}

func removeDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func printAndQuitOnError(err error) {
	if err != nil {
		fmt.Printf("error: %s", err.Error())
		os.Exit(1)
	}
}
