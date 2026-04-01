package torrentFile

import (
	"bytes"
	"crypto/sha1"
	"github.com/jackpal/bencode-go"
	"os"
)

type TorrentFile struct {
	Announce     string
	Comment      string
	CreationDate int
	InfoHash     [20]byte // SHA-1 produces 20-byte hash
	Name         string
	Length       int
	PieceLength  int
	PieceHashes  [][20]byte
}

type BencodeTorrentFileInfo struct {
	Name        string `bencode:"name"`
	Length      int    `bencode:"length"`
	PieceLength int    `bencode:"piece length"`
	Pieces      string `bencode:"pieces"`
}

type BencodeTorrentFile struct {
	Announce     string                 `bencode:"announce"`
	Comment      string                 `bencode:"comment"`
	CreationDate int                    `bencode:"creation date"`
	Info         BencodeTorrentFileInfo `bencode:"info"`
}

func ParseTorrentFile(torrentFilePath string) (TorrentFile, error) {
	/* open file */
	f, err := os.Open(torrentFilePath)
	if err != nil {
		return TorrentFile{}, err
	}
	defer f.Close()

	/* decode bencode */
	btf := BencodeTorrentFile{}
	err = bencode.Unmarshal(f, &btf)
	if err != nil {
		return TorrentFile{}, err
	}

	/* hash the info using SHA-1 */
	buffer := bytes.Buffer{}
	err = bencode.Marshal(&buffer, btf.Info)
	if err != nil {
		return TorrentFile{}, err
	}
	infoHash := sha1.Sum(buffer.Bytes())

	/* split pieces into [20]byte hashes */
	pieces := []byte(btf.Info.Pieces)
	numPieces := len(pieces) / 20
	pieceHashes := make([][20]byte, numPieces)
	for i := 0; i < numPieces; i++ {
		copy(pieceHashes[i][:], pieces[i*20:(i+1)*20])
	}

	torrentFile := TorrentFile{
		Announce:     btf.Announce,
		Comment:      btf.Comment,
		CreationDate: btf.CreationDate,
		InfoHash:     infoHash,
		Name:         btf.Info.Name,
		Length:       btf.Info.Length,
		PieceLength:  btf.Info.PieceLength,
		PieceHashes:  pieceHashes,
	}

	return torrentFile, nil
}
