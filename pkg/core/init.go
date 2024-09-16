package core

import (
	"fmt"
	"os"
	"strings"
)

func Init() {
	if len(os.Args) < 2 || len(os.Args) == 2 && !strings.HasSuffix(os.Args[1], ".txt") {
		fmt.Println("provide txt file")
		os.Exit(1)
	} else if len(os.Args) > 2 {
		fmt.Println("too many arguments")
		os.Exit(1)
	}

	InputHandler(os.Args[1])
}
