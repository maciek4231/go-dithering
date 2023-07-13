package main

import (
	"fmt"
	"math"
	"os"

	vidio "github.com/AlexEidt/Vidio"
)

var bayer = [16][16]int{
	{0, 191, 48, 239, 12, 203, 60, 251, 3, 194, 51, 242, 15, 206, 63, 254},
	{127, 64, 175, 112, 139, 76, 187, 124, 130, 67, 178, 115, 142, 79, 190, 127},
	{32, 223, 16, 207, 44, 235, 28, 219, 35, 226, 19, 210, 47, 238, 31, 222},
	{159, 96, 143, 80, 171, 108, 155, 92, 162, 99, 146, 83, 174, 111, 158, 95},
	{8, 199, 56, 247, 4, 195, 52, 243, 11, 202, 59, 250, 7, 198, 55, 246},
	{135, 72, 183, 120, 131, 68, 179, 116, 138, 75, 186, 123, 134, 71, 182, 119},
	{40, 231, 24, 215, 36, 227, 20, 211, 43, 234, 27, 218, 39, 230, 23, 214},
	{167, 104, 151, 88, 163, 100, 147, 84, 170, 107, 154, 91, 166, 103, 150, 87},
	{2, 193, 50, 241, 14, 205, 62, 253, 1, 192, 49, 240, 13, 204, 61, 252},
	{129, 66, 177, 114, 141, 78, 189, 126, 128, 65, 176, 113, 140, 77, 188, 125},
	{34, 225, 18, 209, 46, 237, 30, 221, 33, 224, 17, 208, 45, 236, 29, 220},
	{161, 98, 145, 82, 173, 110, 157, 94, 160, 97, 144, 81, 172, 109, 156, 93},
	{10, 201, 58, 249, 6, 197, 54, 245, 9, 200, 57, 248, 5, 196, 53, 244},
	{137, 74, 185, 122, 133, 70, 181, 118, 136, 73, 184, 121, 132, 69, 180, 117},
	{42, 233, 26, 217, 38, 229, 22, 213, 41, 232, 25, 216, 37, 228, 21, 212},
	{169, 106, 153, 90, 165, 102, 149, 86, 168, 105, 152, 89, 164, 101, 148, 85}}

func main() {
	switch len(os.Args) {
	case 1:
		fmt.Println("Error: no input/output paths provided")
		os.Exit(1)
	case 2:
		fmt.Println("Error: no output path provided")
		os.Exit(1)
	}
	path := os.Args[1]
	out := os.Args[2]
	video, err := vidio.NewVideo(path)
	if err != nil {
		fmt.Println("error", err)
	}
	options := vidio.Options{
		FPS:     video.FPS(),
		Bitrate: video.Bitrate(),
		Codec:   video.Codec(),
	}
	width := video.Width()
	framesCount := int(video.Duration() * video.FPS())
	writer, err := vidio.NewVideoWriter(
		out,
		width,
		video.Height(),
		&options,
	)

	defer writer.Close()
	leng := width * video.Height() * 4
	frC := 0
	for video.Read() {
		frame := video.FrameBuffer()
		var newFrame []byte
		x := 0
		y := 0
		for i := 0; i < leng; i += 4 {
			percLight := lightness((0.2126*linear(frame[i]) + 0.7152*linear(frame[i+1]) + 0.0722*linear(frame[i+2])))

			if percLight < (float64(bayer[x%16][y%16]) / 256.0 * 100) {
				newFrame = append(newFrame, 0, 0, 0, 255)
			} else {
				newFrame = append(newFrame, 255, 255, 255, 255)
			}

			if (x+1)%width != 0 {
				x++
			} else {
				y++
				x = 0
			}
		}

		frC++
		writer.Write(newFrame)
		fmt.Printf("\033[2K\r%d/%d", frC, framesCount)
	}
	fmt.Println("\nSuccess")
}

func linear(V byte) float64 {
	var L float64 = float64(V) / 255
	if L <= 0.04045 {
		return L / 12.92
	} else {
		return math.Pow((L+0.055)/1.055, 2.4)
	}
}

func lightness(L float64) float64 {
	if L <= 0.008856 {
		return L * 903.3
	} else {
		return math.Pow(L, 1.0/3)*116 - 16
	}
}
