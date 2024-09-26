package main

import (
	"bytes"
	"fmt"
	"io"
	"lem-in/pkg/visualizer"
	"os"
)

func main() {
	// Lire l'entrée complète
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Printf("Erreur lors de la lecture de l'entrée: %v\n", err)
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
