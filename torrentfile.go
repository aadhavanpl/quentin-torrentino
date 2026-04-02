package torrentFile

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/jackpal/bencode-go"
	"github.com/veggiedefender/torrent-client/peers"
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

type BencodeTrackerResponse struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

type Peer struct {
    IP   net.IP // 4 butes
    Port uint16 // 2 bytes - Big-endian
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

func BuildTrackerUrl(torrentFile TorrentFile) (string, error) {
	base, err := url.Parse(torrentFile.Announce)
    if err != nil {
        return "", err
    }

	params := url.Values{
        "info_hash":  []string{string(torrentFile.InfoHash[:])},
        "peer_id":    []string{string(peerID[:])},
        "port":       []string{strconv.Itoa(int(Port))},
        "uploaded":   []string{"0"},
        "downloaded": []string{"0"},
        "compact":    []string{"1"},
        "left":       []string{strconv.Itoa(torrentFile.Length)},
    }

	base.RawQuery = params.Encode()
    return base.String(), nil
}

func RequestPeers(trackerUrl string) ([]Peer, error) {
	httpClient := &http.Client{Timeout: 15 * time.Second}
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	trackerResponse := BencodeTrackerResponse{}
	err = bencode.Unmarshal(resp.Body, &trackerResponse)
	if err != nil {
		return nil, err
	}

	return peers.Unmarshal([]byte(trackerResponse.Peers))
}

func ParseTrackerResponse(peersData []byte) ([]Peer, error) {
	const peerSize = 6

	if len(peersData)%peerSize != 0 {
        return nil, fmt.Errorf("Received malformed peers")
    }

	numPeers = len(peersData) / peerSize
	peers := make([]Peer, numPeers)
	for i := 0; i < numPeers; i++ {
        offset := i * peerSize
        peers[i].IP = net.IP(peersData[offset : offset+4])
        peers[i].Port = binary.BigEndian.Uint16(peersData[offset+4 : offset+6])
    }
    return peers, nil
}