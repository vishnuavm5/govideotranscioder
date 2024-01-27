package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/joho/godotenv"
)

type Resolution struct {
	Height string
	Width  string
}

func convertVideo(inputPath, outputPath string, height, width string) {
	ffmpegCommand := fmt.Sprintf("ffmpeg -i %s -y -acodec aac -vcodec libx264 -filter:v scale=w=%s:h=%s -f mp4 %s", inputPath, width, height, outputPath)

	cmd := exec.Command("bash", "-c", ffmpegCommand)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error creating StdoutPipe:", err)
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println("Error creating StderrPipe:", err)
		return
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting command:", err)
		return
	}

	go func() {
		for {
			buffer := make([]byte, 1024)
			n, err := stdout.Read(buffer)
			if n > 0 {
				fmt.Printf("stdout: %s", buffer[:n])
			}
			if err != nil {
				break
			}
		}
	}()

	go func() {
		for {
			buffer := make([]byte, 1024)
			n, err := stderr.Read(buffer)
			if n > 0 {
				fmt.Printf("stderr: %s", buffer[:n])
			}
			if err != nil {
				break
			}
		}
	}()

	if err := cmd.Wait(); err != nil {
		fmt.Printf("Command finished with error: %v\n", err)
	}
}

func main() {
	godotenv.Load()
	inputVideoPath := os.Getenv("INPUT_DIR")
	outputDir := os.Getenv("OUTPUT_DIR")

	resolutions := []Resolution{
		{Height: "1080", Width: "1920"},
		{Height: "720", Width: "1280"},
		{Height: "480", Width: "640"},
		{Height: "360", Width: "480"},
	}

	convertVideo(inputVideoPath, fmt.Sprintf("%svideo.mp4", outputDir), "1080", "1920")

	for _, resolution := range resolutions {
		outputVideoPath := fmt.Sprintf("%svideo_%sp.mp4", outputDir, resolution.Height)
		convertVideo(inputVideoPath, outputVideoPath, resolution.Height, resolution.Width)
	}
}
