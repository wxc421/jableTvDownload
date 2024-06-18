package main

import (
	"github.com/wxc421/jableTvDownload/client"
	"github.com/wxc421/jableTvDownload/m3u8/download"
	"github.com/wxc421/jableTvDownload/m3u8/parse"
	"log"
	"path/filepath"
)

func main() {

	downloadUrl := "xxxxxx" // input jableTv url here
	proxyUrl := "http://127.0.0.1:7890"
	saveDir := "D:\\movies"
	concurrency := 24

	client.SetProxy(proxyUrl)

	c, err := client.GetProxyClient()
	if err != nil {
		panic(err)
	}
	// get m3u8 url and title
	resp, err := c.R().Get(downloadUrl)
	if err != nil {
		panic(err)
	}
	m3u8Url := parse.FindM3u8(resp.Body())
	title := parse.FindTitle(resp.Body())
	log.Printf("m3u8Url:%v", m3u8Url)
	log.Printf("title:%v", title)
	savePath := filepath.Join(saveDir, title)
	if len(m3u8Url) > 0 {
		task, err := download.NewTask(savePath, m3u8Url)
		if err != nil {
			log.Fatal(err)
			return
		}
		if err := task.Start(concurrency); err != nil {
			log.Fatal(err)
			return
		}
		log.Printf("[%v] download success:%v", title, savePath)
	}
}
