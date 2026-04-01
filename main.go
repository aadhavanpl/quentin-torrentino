package main

import (
	"github.com/aadhavanpl/quentin-torrentino/torrentfile"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Torrent file path required!")
	}
	torrentFilePath := os.Args[1]

	torrentFile, err := torrentfile.ParseTorrentFile(torrentFilePath)
	if err != nil {
		log.Fatal(err)
	}
}
