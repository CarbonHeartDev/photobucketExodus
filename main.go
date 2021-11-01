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

	workDir, err := os.Getwd()
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

	downloadResults := make(map[string]DownloadResult)
	for i, match := range matches {
		fmt.Printf("Downloading file %d of %d: %s\n", i+1, len(matches), match)
		downloadResults[match] = downloadAndWriteFile(match, outDir)
	}

	fmt.Printf("Data saved in %s\n", outDir)
	os.WriteFile(filepath.Join(workDir, "downloadReport.csv"), []byte(writeCsvReport(downloadResults)), 0664)
	fmt.Printf("A report of the downloads was written in %s\n\n", filepath.Join(workDir))

}

func downloadAndWriteFile(url, outDir string) DownloadResult {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "image/webp,*/*")
	req.Header.Add("Referer", url)

	resp, err := client.Do(req)
	if err != nil {
		return DownloadResult{url, "", "Error in processing the HTTP request " + err.Error()}
	}

	if resp.StatusCode != 200 {
		errorMsg := "Download failed HTTP code " + resp.Status
		fmt.Println(errorMsg)
		return DownloadResult{url, "", errorMsg}
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return DownloadResult{url, "", "Error in processing the response data " + err.Error()}
	}

	urlSplits := strings.Split(url, "/")
	fileName := urlSplits[len(urlSplits)-1]
	fileExtension := strings.Split(fileName, ".")[len(strings.Split(fileName, "."))-1]

	if fileExtension == "webp" {
		err = os.WriteFile(filepath.Join(outDir, fileName), respBody, 0600)
		if err != nil {
			return DownloadResult{url, "", "Error in saving the file " + err.Error()}
		}

		return DownloadResult{url, fileName, ""}

	} else {
		format, err := imgconv.FormatFromExtension(fileExtension)
		if err != nil {
			return DownloadResult{url, "", "Error in detecting image extension from name"}
		}

		buf := new(bytes.Buffer)
		img, err := imgconv.Decode(bytes.NewReader(respBody))
		err = imgconv.Write(buf, img, imgconv.FormatOption{Format: format})
		if err != nil {
			return DownloadResult{url, "", "Error in converting the image"}
		}
		err = os.WriteFile(filepath.Join(outDir, fileName), buf.Bytes(), 0600)
		if err != nil {
			return DownloadResult{url, "", "Error in writing the output file"}
		}

		return DownloadResult{url, fileName, ""}

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

func writeCsvReport(downloadResults map[string]DownloadResult) string {

	var sb strings.Builder

	sb.WriteString("url,outputFileName,error\n")
	for _, downloadResult := range downloadResults {
		sb.WriteString(downloadResult.url + "," + downloadResult.outputFileName + "," + downloadResult.error + "\n")
	}

	return sb.String()
}

type DownloadResult struct {
	url, outputFileName, error string
}
