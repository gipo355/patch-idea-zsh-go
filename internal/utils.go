package internal

import (
	"bufio"
	"os"
	"strings"
)

func ReadLine() string {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}
