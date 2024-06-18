package parse

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/wxc421/jableTvDownload/client"
	"io/ioutil"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/oopsguy/m3u8/tool"
)

type Result struct {
	URL  *url.URL
	M3u8 *M3u8
	Keys map[int]string
}

func FindM3u8(content []byte) string {
	groups := regexp.MustCompile(`http[s]://[a-zA-Z0-9/\\.%_-]+.m3u8`).FindSubmatch(content)
	return string(groups[0])
}

func FindTitle(content []byte) string {
	compileRegex := regexp.MustCompile("<title>(.*?) - Jable.TV.*</title>")
	matchArr := compileRegex.FindStringSubmatch(string(content))

	if len(matchArr) > 0 {
		return strings.ReplaceAll(matchArr[len(matchArr)-1], " ", "_")
	}
	return ""
}

func ParseM3u8FromUrl(link string) (*Result, error) {
	u, err := url.Parse(link)
	if err != nil {
		return nil, err
	}
	link = u.String()
	c, _ := client.GetClient()
	resp, err := c.R().Get(link)
	if err != nil {
		return nil, fmt.Errorf("request m3u8 URL failed: %s", err.Error())
	}
	data := resp.Body()
	// for test
	_ = os.WriteFile("result.m3u8", data, 0644)
	//noinspection GoUnhandledErrorResult
	m3u8, err := parse(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	if len(m3u8.MasterPlaylist) != 0 {
		sf := m3u8.MasterPlaylist[0]
		return ParseM3u8FromUrl(tool.ResolveURL(u, sf.URI))
	}
	if len(m3u8.Segments) == 0 {
		return nil, errors.New("can not found any TS file description")
	}
	result := &Result{
		URL:  u,
		M3u8: m3u8,
		Keys: make(map[int]string),
	}

	for idx, key := range m3u8.Keys {
		switch {
		case key.Method == "" || key.Method == CryptMethodNONE:
			continue
		case key.Method == CryptMethodAES:
			// Request URL to extract decryption key
			keyURL := key.URI
			keyURL = tool.ResolveURL(u, keyURL)
			resp, err := tool.Get(keyURL)
			if err != nil {
				return nil, fmt.Errorf("extract key failed: %s", err.Error())
			}
			keyByte, err := ioutil.ReadAll(resp)
			_ = resp.Close()
			if err != nil {
				return nil, err
			}
			// fmt.Println("decryption key: ", string(keyByte))
			result.Keys[idx] = string(keyByte)
		default:
			return nil, fmt.Errorf("unknown or unsupported cryption method: %s", key.Method)
		}
	}
	return result, nil
}
