package main

// go build -ldflags -H=windowsgui -o screentime.exe  main.go

import (
	"bytes"
	"github.com/shirou/gopsutil/process"
	"github.com/gen2brain/beeep"
	"image/jpeg"
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

func isAlreadyRunning() bool {
	currentProcess, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		return false
	}
	
	currentName, err := currentProcess.Name()
	if err != nil {
		return false
	}
	
	processes, err := process.Processes()
	if err != nil {
		return false
	}

	for _, p := range processes {
		if p.Pid == currentProcess.Pid {
			continue
		}
		name, err := p.Name()
		if err == nil && name == currentName {
			return true
		}
	}
	return false
}

func showNotification(title, message string) {
	beeep.Notify(title, message, "")
}

func captureScreenshot(filename string) {
	n := screenshot.NumActiveDisplays()

	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)

		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			log.Fatalf("Error while capturing screenshot: %v", err)
		}
		file, err := os.Create(filename)
		if err != nil {
			log.Fatalf("Error creating file: %v", err)
		}
		defer file.Close()

		err = jpeg.Encode(file, img, &jpeg.Options{Quality: 18})
		if err != nil {
			log.Printf("Error encoding jpeg: %v", err)
			return
		}

		log.Printf("#%d : %v \"%s\"\n", i, bounds, filename)
	}
}

func uploadFile(url string, filename string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("Error reading file: %v", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		log.Printf("Error creating HTTP request: %v", err)
		consecutiveFailures++
		return
	}
	req.Header.Set("Content-Type", "multipart/form-data")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		//log.Printf("Error doing HTTP request: %v", err)
		consecutiveFailures++
		return
	} else {
		consecutiveFailures = 0
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return
	}

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
	if isAlreadyRunning() {
		// Silently exit if another instance is running
		showNotification("info", "screentime already running")
		return
	}

	showNotification("info", "screentime is running, please focus on your study")

	appdata := os.Getenv("APPDATA")
	dir := filepath.Join(appdata, appName)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Fatalf("Error creating directory: %v", err)
	}

	t := time.Now()
	datePath := filepath.Join(dir, t.Format("2006-01-02"))
	if err := os.MkdirAll(datePath, os.ModePerm); err != nil {
		log.Fatalf("Error creating date directory: %v", err)
	}

	logFile, err := os.OpenFile(filepath.Join(datePath, "app.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer logFile.Close()

	// Only log to the file
	log.SetOutput(logFile)

	// Delete screenshots older than 3 months
	threshold := t.AddDate(0, 0, -90)

	deleteOldScreenshots(dir, threshold)

	for {
		t := time.Now()
		fileName := filepath.Join(datePath, "screenshot_"+t.Format("15-04-05")+".jpg")

		captureScreenshot(fileName)
		if consecutiveFailures < maxRetryCount {
			uploadFile("http://192.168.0.21:5000/upload", fileName)
		} else {
			//log.Printf("Reached maximum number of consecutive upload failures. Skipping uploads.")
		}

		time.Sleep(20 * time.Second)
	}
}
