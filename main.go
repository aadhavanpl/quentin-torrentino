package main

import (
	"log"
	"os"

	"github.com/aadhavanpl/quentin-torrentino/torrentfile"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Torrent file path required!")
	}
	torrentFilePath := os.Args[1]
	outputPath := os.Args[2]

	torrentFile, err := torrentfile.ParseTorrentFile(torrentFilePath)
	if err != nil {
		log.Fatal(err)
	}

	trackerUrl, err := torrentFile.BuildTrackerUrl(torrentFile)
	if err != nil {
		log.Fatal(err)
	}

	peers, err := torrentFile.RequestPeers(trackerUrl)
	if err != nil {
		log.Fatal(err)
	}
}
