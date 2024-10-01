## lem-in

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/Shasor/lem-in)

### Description

This project is meant to make we code a digital version of an ant farm.
Lem-in will read from a file (describing the ants and the colony) given in the arguments.
Upon successfully finding the quickest path, lem-in will display the content of the file passed as argument and each move the ants make from room to room.

visualizer reads Lem-in's output, then launches a graphical display of the anthill and its resolution.

### Usage

```bash
$ make # to build
Build in progress...
OK!
$ ./lem-in <file_name>.txt
...
$ ./lem-in <file_name>.txt | ./visualizer
...
```

### Author(s)

- [Théo VALLOIS](https://zone01normandie.org/git/tvallois)
- [Arnaud de Dreuzy](https://zone01normandie.org/git/adedreuz)
- [Adam GONÇALVES](https://zone01normandie.org/git/agoncalv) (aka [Shasor#3755](https://discordapp.com/users/282816260075683841))
