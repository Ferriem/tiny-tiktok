package ffmpeg

import (
	"bytes"
	"os"
	"os/exec"
)

func Cover(videoUrl string, timeOffest string) ([]byte, error) {
	cmd := exec.Command(
		"ffmpeg", "-i", videoUrl, "-ss", timeOffest, "-vframes", "1", "-q:v", "2", "-f", "image2", "pipe:1",
	)
	cmd.Stderr = os.Stderr

	var outputBuffer bytes.Buffer
	cmd.Stdout = &outputBuffer

	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return outputBuffer.Bytes(), nil
}
