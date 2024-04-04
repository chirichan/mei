package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/chirichan/rice"
)

const BaseUrl = "https://cn.bing.com"

var (
	urls = []string{
		"https://cn.bing.com/HPImageArchive.aspx?format=js&idx=0&n=7",
		"https://cn.bing.com/HPImageArchive.aspx?format=js&idx=8&n=8",
	}
	outputDir = flag.String("o", "", "Output directory, default output: \"~/Pictures/Saved Pictures/BingImages\"")
)

func init() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
}

type Image struct {
	Startdate     string `json:"startdate"`
	Fullstartdate string `json:"fullstartdate"`
	Enddate       string `json:"enddate"`
	Url           string `json:"url"`
	Urlbase       string `json:"urlbase"`
	Copyright     string `json:"copyright"`
	Copyrightlink string `json:"copyrightlink"`
	Title         string `json:"title"`
	Quiz          string `json:"quiz"`
	Hsh           string `json:"hsh"`
}

type HPImageArchive struct {
	Images   []Image `json:"images"`
	Tooltips struct {
		Loading  string `json:"loading"`
		Previous string `json:"previous"`
		Next     string `json:"next"`
		Walle    string `json:"walle"`
		Walls    string `json:"walls"`
	} `json:"tooltips"`
}

func main() {
	flag.Parse()
	if *outputDir == "" {
		_ = os.MkdirAll("BingImages", 0644)
		*outputDir = "BingImages"
	}
	log.Printf("output dir: %s\n", *outputDir)
	if err := os.MkdirAll(*outputDir, 0644); err != nil {
		log.Fatalf("create dir %s err: %+v", *outputDir, err)
	}

	var wg = &sync.WaitGroup{}
	var imageChan = make(chan Image)
	for _, v := range urls {
		wg.Add(1)
		go requestBingAPI(wg, v, imageChan)
	}
	go func() {
		wg.Wait()
		close(imageChan)
		log.Println("image chan closed")
	}()
	var images []Image
	for res := range imageChan {
		images = append(images, res)
	}
	saveImageJson(images)

	for _, image := range images {
		wg.Add(1)
		go downloadAndSaveBingImage(wg, image)
	}
	wg.Wait()
}

func saveImageJson(_ []Image) error {
	// TODO save json by year
	return nil
}

func requestBingAPI(wg *sync.WaitGroup, url string, imageChan chan Image) {
	defer wg.Done()

	body, err := getBody(url)
	if err != nil {
		log.Printf("get body %s err: %v", url, err)
		return
	}
	var hpa HPImageArchive
	if err := json.Unmarshal(body, &hpa); err != nil {
		log.Printf("json unmarshal body into hpa err: %v", err)
		return
	}

	for _, v := range hpa.Images {
		imageChan <- v
	}
}

func downloadAndSaveBingImage(wg *sync.WaitGroup, image Image) {
	defer wg.Done()
	hdUrl := BaseUrl + image.Urlbase + "_1920x1080.jpg"
	uhdUrl := BaseUrl + image.Urlbase + "_UHD.jpg"
	hd, err := url.Parse(hdUrl)
	if err != nil {
		log.Printf("parse url %s err: %+v", hdUrl, err)
		return
	}
	uhd, err := url.Parse(uhdUrl)
	if err != nil {
		log.Printf("parse url %s err: %+v", uhdUrl, err)
		return
	}
	filename1 := filepath.Join(*outputDir, strings.TrimPrefix(hd.Query().Get("id"), "OHR."))
	filename2 := filepath.Join(*outputDir, strings.TrimPrefix(uhd.Query().Get("id"), "OHR."))
	writeFunc := func(filename, imageUrl string) error {
		if rice.FileExists(filename) {
			return nil
		}
		begin := time.Now()
		log.Printf("⏳ start save %s\n", filename)
		if fileExists(filename) {
			log.Printf("🧐 already exists %s\n", filename)
			return nil
		}
		b, err := getBody(imageUrl)
		if err != nil {
			return fmt.Errorf("get body %s err: %w", imageUrl, err)
		}
		if err := write(filename, b); err != nil {
			return fmt.Errorf("save %s err: %w", filename, err)
		}
		log.Printf("🎉 saved success %s, cost %s\n", filename, time.Since(begin))
		return nil
	}
	var wg2 = &sync.WaitGroup{}
	wg2.Add(2)
	go func(wg2 *sync.WaitGroup, filename, url2 string) {
		defer wg2.Done()
		if err := writeFunc(filename, url2); err != nil {
			log.Printf("err: %+v", err)
		}
	}(wg2, filename1, hdUrl)
	go func(wg2 *sync.WaitGroup, filename, url2 string) {
		defer wg2.Done()
		if err := writeFunc(filename, url2); err != nil {
			log.Printf("err: %+v", err)
		}
	}(wg2, filename2, uhdUrl)
	wg2.Wait()
}

func write(filename string, b []byte) error {
	if fileExists(filename) {
		return nil
	}
	if err := os.WriteFile(filename, b, 0644); err != nil {
		return fmt.Errorf("write file %s err: %w", filename, err)
	}
	return nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func getBody(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	readAll, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return readAll, nil
}

func getHomeDir() string {
	userHomeDir, _ := os.UserHomeDir()
	return userHomeDir
}
