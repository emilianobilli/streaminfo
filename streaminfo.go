package streaminfo

import (
	"encoding/json"
	"fmt"
)

type Info struct {
	value map[string]interface{}
}

type InfoInterface interface {
	GetKeys() []string
	Get(string) (any, bool)
}

type StreamInfo struct {
	Format InfoInterface
	Video  InfoInterface
	Audio  InfoInterface
}

func (s StreamInfo) HasVideo() bool {
	return s.Video != nil
}

func (s StreamInfo) HasAudio() bool {
	return s.Audio != nil
}

func (i *Info) GetKeys() []string {
	keys := []string{}
	if i != nil {
		for k, _ := range i.value {
			keys = append(keys, k)
		}
	}
	return keys
}

func (i *Info) Get(k string) (any, bool) {
	if i != nil {
		if value, ok := i.value[k]; ok {
			return value, true
		}
	}
	return nil, false
}

func extractInfo(streams []interface{}, ctype string) InfoInterface {
	for _, s := range streams {
		if v, ok := s.(map[string]interface{}); ok {
			if codec_type, ok := v["codec_type"]; ok && ctype == codec_type {
				return &Info{v}
			}
		}
	}
	return nil
}

func ExtractStreamInfo(filename string) (*StreamInfo, error) {
	buf, err := FFPROBE(filename)
	if err != nil {
		return nil, err
	}
	rawInfo := make(map[string]interface{})
	if e := json.Unmarshal(buf.Bytes(), &rawInfo); e != nil {
		return nil, e
	}

	var info StreamInfo

	if streams, ok := rawInfo["streams"]; ok {

		if s, ok := streams.([]interface{}); ok {
			info.Video = extractInfo(s, "video")
			info.Audio = extractInfo(s, "audio")
		}
	}

	if info.Video == nil {
		return nil, fmt.Errorf("wrong video information")
	}

	if info.Audio == nil {
		return nil, fmt.Errorf("wrong audio information")
	}

	if format, ok := rawInfo["format"]; ok {
		info.Format = &Info{format.(map[string]interface{})}
	}
	return &info, nil
}
