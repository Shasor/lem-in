package main

import (
	"bytes"
	"fmt"
	"io"
	"lem-in/pkg/visualizer"
	"os"
	"strconv"
	"strings"
)

func main() {
	// exit if no pipe
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		fmt.Println("Aucune entrée disponible sur stdin. Utilisez un pipe pour fournir une entrée.")
		os.Exit(1)
	}

	// read Stdin
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Printf("Erreur lors de la lecture de l'entrée: %v\n", err)
		os.Exit(1)
	}
	// test if lem-in's output isn't already an error.
	if strings.HasPrefix(string(input), "ERROR") {
		fmt.Print(string(input))
		os.Exit(1)
	}
	if _, err := strconv.Atoi(string(input[0])); err != nil {
		fmt.Println("ERROR: invalid data format")
		os.Exit(1)
	}

	fmt.Print(string(input))

	// Créer un nouveau lecteur avec l'entrée lue
	reader := bytes.NewReader(input)
	// Créer une fonction qui renvoie le lecteur
	getReader := func() io.Reader {
		_, err := reader.Seek(0, io.SeekStart)
		if err != nil {
			fmt.Printf("Erreur lors de la réinitialisation du lecteur: %v\n", err)
			os.Exit(1)
		}
		// fmt.Print(string(input))
		return reader
	}
	// Passer la fonction getReader à NewApp
	app, err := visualizer.NewApp(800, 600, getReader)
	if err != nil {
		fmt.Printf("Erreur lors de l'initialisation de l'application : %v\n", err)
		os.Exit(1)
	}
	if err := app.Run(); err != nil {
		fmt.Printf("Erreur lors de l'exécution de l'application : %v\n", err)
		os.Exit(1)
	}
}
