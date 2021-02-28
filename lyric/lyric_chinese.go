package lyric

import (
	// "encoding/json"
	"fmt"
	// "net/url"
	"io"
	"os"
	"strconv"

	"github.com/asmcos/requests"
	"github.com/martinlindhe/subtitles"
	// "github.com/gocolly/colly"
	// r "github.com/solos/requests"
)

// GetLyricOptionsChinese queries available song lyrics. It returns map of title and
// id of the lyric.
func GetLyricOptionsChinese(search string, serviceProvider string) (map[string]string, error) {

	result := make(map[string]string)
	p := requests.Params{
		"site":   serviceProvider,
		"search": search,
	}
	req := requests.Requests()
	req.Header.Set("Content-Type", "application/json")
	resp, err := req.Get("http://api.sunyj.xyz", p)
	if err != nil {
		return nil, err
	}

	var dataMap []map[string]interface{}
	err = resp.Json(&dataMap)
	if err != nil {
		return nil, err
	}
	for _, v := range dataMap {
		songName := v["name"]
		songArtist := v["artist"]
		var lyricID string
		if serviceProvider == "netease" {
			lyricIDfloat64 := v["lyric_id"]
			lyricID = strconv.FormatFloat(lyricIDfloat64.(float64), 'f', -1, 64)
		} else if serviceProvider == "kugou" {
			lyricID = v["lyric_id"].(string)
		}
		songTitle := fmt.Sprintf("%s - %s ", songArtist, songName)
		if lyricID == "" {
			continue
		}
		result[songTitle] = lyricID
	}

	return result, nil
}

// GetLyricChinese should receive url that was returned from GetLyricOptions. GetLyric
// returns lyric of the queried song.
func GetLyricChinese(lyricID string, serviceProvider string) (string, error) {

	var lyric string
	p := requests.Params{
		"site":  serviceProvider,
		"lyric": lyricID,
	}
	req := requests.Requests()
	resp, err := req.Get("http://api.sunyj.xyz", p)
	if err != nil {
		return "", err
	}
	var dataMap map[string]interface{}
	err = resp.Json(&dataMap)
	if err != nil {
		return "", err
	}
	lyric = dataMap["lyric"].(string)
	if looksLikeLRC(lyric) {
		//Fixme
		filename := "/home/tramhao/1.lrc"
		file, _ := os.Create(filename)
		io.WriteString(file, lyric)
		file.Close()

		var tmpSubtitle subtitles.Subtitle
		tmpSubtitle, err = NewFromLRC(lyric)
		if err != nil {
			return "", err
		}
		lyric = tmpSubtitle.AsSRT()
		//Fixme
		filename = "/home/tramhao/1.srt"
		file, _ = os.Create(filename)
		io.WriteString(file, lyric)
		file.Close()

	}
	return lyric, nil
}
