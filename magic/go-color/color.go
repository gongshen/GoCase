package main

import "fmt"

const (
	red = uint8(91 + iota)
	green
)

func main() {
	Print(green)
}

func Print(color uint8) {
	fmt.Printf("\x1b[%dmI ♡  You \x1b[0m", color)
}
