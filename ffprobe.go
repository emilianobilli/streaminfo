package streaminfo

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

var ffprobe_exec string

func init() {
	path, _ := os.Getwd()
	if runtime.GOOS == "windows" {
		ffprobe_exec = path + "/ffprobe.exe"
	} else {
		ffprobe_exec = path + "/ffprobe"
	}
}

func ExistBinaryFile() bool {
	info, err := os.Stat(ffprobe_exec)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func FFPROBE(filename string) (*bytes.Buffer, error) {
	if !ExistBinaryFile() {
		return nil, fmt.Errorf("ffprobe is not installed")
	}
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("filename: %s does not exist", filename)
	}

	buf := &bytes.Buffer{}

	cmd := exec.Command(ffprobe_exec, "-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", "-i", filename)
	cmd.Stdout = buf
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return buf, nil
}
