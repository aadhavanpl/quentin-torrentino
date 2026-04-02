package main

import (
	"log"
	"os"

	"github.com/aadhavanpl/quentin-torrentino/torrent"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Torrent file path required!")
	}
	torrentFilePath := os.Args[1]
	outputPath := os.Args[2]

	torrentFile, err := torrent.ParseTorrentFile(torrentFilePath)
	if err != nil {
		log.Fatal(err)
	}

	peerId, err := torrent.GeneratePeerId()
	if err != nil {
		log.Fatal(err)
	}

	trackerUrl, err := torrent.BuildTrackerUrl(torrentFile, peerId)
	if err != nil {
		log.Fatal(err)
	}

	peers, err := torrent.RequestPeers(trackerUrl)
	if err != nil {
		log.Fatal(err)
	}

	err = torrent.DownloadFile(torrentFile, peers, peerId, outputPath)
	if err != nil {
		log.Fatal(err)
	}
}
