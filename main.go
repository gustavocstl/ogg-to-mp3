package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"
)

func oggFiles() []string {
	fileDivider := "%DV%"

	cmd := "ls | grep -E '\\.ogg(\\n|$)'"

	out, err := exec.Command("bash", "-c", cmd).Output()

	if err != nil {
		panic(err.Error())
	}

	filesBuf := bytes.Buffer{}

	for i := 0; i < len(out); i++ {
		if out[i] == 10 {
			filesBuf.Write([]byte(fileDivider))
			continue
		}

		filesBuf.Write([]byte{out[i]})
	}

	files := strings.Split(filesBuf.String(), fileDivider)

	return files
}

func convert(fileName string, folderName string, wg *sync.WaitGroup) {
	if fileName == "" {
		wg.Done()
		return
	}

	fmt.Printf("Converting %s to mp3...\n", fileName)

	cmd := "ffmpeg -i " + fileName + " -f mp3 " + folderName + "/" + fileName + ".mp3"

	out, err := exec.Command("bash", "-c", cmd).Output()

	if err != nil {
		fmt.Println(out, cmd)
		panic(err)
	}

	wg.Done()
}

func main() {
	var wg sync.WaitGroup
	files := oggFiles()
	startTime := time.Now()
	folderName := "mp3-files-" + startTime.Format("20060102150405")

	exec.Command("mkdir", folderName).Run()

	fmt.Printf(
		"Starting conversion .ogg to .mp3 for %d files - [%s]\n\n",
		len(files),
		startTime.Format("02/01/2006 15:04:05"),
	)

	for _, file := range files {
		wg.Add(1)
		go convert(file, folderName, &wg)
	}

	wg.Wait()

	fmt.Printf("All files are formatted and saved in "+folderName+" folder.\nTime to format: %f seconds\n", time.Since(startTime).Seconds())
}
