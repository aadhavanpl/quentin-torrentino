package torrent

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/jackpal/bencode-go"
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

type Torrent struct {
	Peers       []Peer
	PeerID      [20]byte
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
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

const PORT = 6000

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

func GeneratePeerId() ([20]byte, error) {
	var peerId [20]byte
	
	_, err := rand.Read(peerId[:])
	if err != nil {
		return [20]byte{}, err
	}

	return peerId, nil
}

func BuildTrackerUrl(torrentFile TorrentFile, peerId [20]byte) (string, error) {
	base, err := url.Parse(torrentFile.Announce)
    if err != nil {
        return "", err
    }

	params := url.Values{
        "info_hash":  []string{string(torrentFile.InfoHash[:])},
        "peer_id":    []string{string(peerId[:])},
        "port":       []string{strconv.Itoa(int(PORT))},
        "uploaded":   []string{"0"},
        "downloaded": []string{"0"},
        "compact":    []string{"1"},
        "left":       []string{strconv.Itoa(torrentFile.Length)},
    }

	base.RawQuery = params.Encode()
    return base.String(), nil
}

func RequestPeers(trackerUrl string) ([]Peer, error) {
	httpClient := &http.Client{}
	resp, err := httpClient.Get(trackerUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	trackerResponse := BencodeTrackerResponse{}
	err = bencode.Unmarshal(resp.Body, &trackerResponse)
	if err != nil {
		return nil, err
	}

	return ParseTrackerResponse([]byte(trackerResponse.Peers))
}

func ParseTrackerResponse(peersData []byte) ([]Peer, error) {
	const peerSize = 6

	if len(peersData)%peerSize != 0 {
        return nil, fmt.Errorf("Received malformed peers")
    }

	numPeers := len(peersData) / peerSize
	peers := make([]Peer, numPeers)
	for i := range numPeers {
        offset := i * peerSize
        peers[i].IP = net.IP(peersData[offset : offset+4])
        peers[i].Port = binary.BigEndian.Uint16(peersData[offset+4 : offset+6])
    }
    return peers, nil
}

func DownloadFile(torrentFile TorrentFile, peers []Peer, peerId [20]byte, outputPath string) error {
	torrent := Torrent {
		Peers: 		 peers,
		PeerID: 	 peerId,
		InfoHash:    torrentFile.InfoHash,
		PieceHashes: torrentFile.PieceHashes,
		PieceLength: torrentFile.PieceLength,
		Length:      torrentFile.Length,
		Name:        torrentFile.Name,
	}
	

}