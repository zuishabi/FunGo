package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func isVideo(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".mp4", ".mkv", ".avi", ".mov", ".flv", ".wmv", ".ts":
		return true
	default:
		return false
	}
}

func main() {
	ffmpegPath := `C:\projects\ffmpeg-master-latest-win64-gpl-shared\bin\ffmpeg.exe`
	inputPath := `C:\Users\张俏悦\Downloads\1\[DMG&Sakurato][Yofukashi_no_Uta_S2][01-12][1080P][GB][MP4]`

	err := filepath.WalkDir(inputPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !isVideo(path) {
			return nil
		}

		parentDir := filepath.Dir(path)
		base := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		outDir := filepath.Join(parentDir, base)

		if mkErr := os.MkdirAll(outDir, 0755); mkErr != nil {
			return mkErr
		}

		outM3u8 := filepath.Join(outDir, "index.m3u8")

		cmd := exec.Command(
			ffmpegPath,
			"-y",
			"-i", path,
			"-c:v", "libx264",
			"-c:a", "aac",
			"-hls_time", "6",
			"-hls_playlist_type", "vod",
			"-hls_segment_filename", filepath.Join(outDir, "segment_%03d.ts"),
			outM3u8,
		)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		fmt.Println("Processing:", path)
		return cmd.Run()
	})

	if err != nil {
		fmt.Println("Error:", err)
	}
}
