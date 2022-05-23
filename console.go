package dango

import (
	"bufio"
	"fmt"
	"os"
)

// ConsolePrompt takes a text prompt, wait for and return user input
func ConsolePrompt(text string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	fmt.Println(text)
	s, _ = r.ReadString('\n')
	return s
}
