package input

import (
	"fmt"
	"os"
	"strings"
)

func InputInit() {
	// gestion arguments + extension .txt
	if len(os.Args) < 2 || len(os.Args) == 2 && !strings.HasSuffix(os.Args[1], ".txt") {
		fmt.Println("provide txt file")
		os.Exit(1)
	} else if len(os.Args) > 2 {
		fmt.Println("too many arguments")
		os.Exit(1)
	}
	// création des rooms et des ants --> voir input.go
	InputHandler(os.Args[1])
}
