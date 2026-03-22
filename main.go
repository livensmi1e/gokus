package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	width, height := 20, 10
	x, y := 0, 0
	dx, dy := 1, 1

	// Bật alternate buffer + ẩn cursor
	fmt.Print("\033[?1049h")
	fmt.Print("\033[?25l")

	// Khi exit (Ctrl+C), bật lại cursor + tắt alternate buffer
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Print("\033[?25h")   // hiện cursor
		fmt.Print("\033[?1049l") // tắt alternate buffer
		os.Exit(0)
	}()

	for {
		frame := "┌" + repeat("─", width) + "┐\n"
		for i := 0; i < height; i++ {
			frame += "│"
			for j := 0; j < width; j++ {
				if i == y && j == x {
					frame += "*"
				} else {
					frame += " "
				}
			}
			frame += "│\n"
		}
		frame += "└" + repeat("─", width) + "┘\n"

		fmt.Print("\033[H") // move cursor top-left
		fmt.Print(frame)

		// update position
		x += dx
		y += dy
		if x <= 0 || x >= width-1 {
			dx = -dx
		}
		if y <= 0 || y >= height-1 {
			dy = -dy
		}

		time.Sleep(200 * time.Millisecond)
	}
}

func repeat(s string, n int) string {
	res := ""
	for i := 0; i < n; i++ {
		res += s
	}
	return res
}
