package torrent

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

// Client holds the connection and state for a single peer
type Client struct {
	Conn     net.Conn
	Choked   bool
	Bitfield Bitfield
	peer     Peer
	infoHash [20]byte
	peerID   [20]byte
}

// NewClient dials the peer, completes the handshake and reads the bitfield
func NewClient(peer Peer, peerID, infoHash [20]byte) (*Client, error) {
	conn, err := net.DialTimeout("tcp", peer.String(), 3*time.Second)
	if err != nil {
		return nil, err
	}

	if _, err := CompleteHandshake(conn, infoHash, peerID); err != nil {
		conn.Close()
		return nil, err
	}

	bf, err := recvBitfield(conn)
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &Client{
		Conn:     conn,
		Choked:   true,
		Bitfield: bf,
		peer:     peer,
		infoHash: infoHash,
		peerID:   peerID,
	}, nil
}

// Read reads a message from the peer
func (c *Client) Read() (*Message, error) {
	return ReadMessage(c.Conn)
}

// SendRequest sends a request for a block
func (c *Client) SendRequest(index, begin, length int) error {
	req := ParsePieceRequest(index, begin, length)
	_, err := c.Conn.Write(req.Serialize())
	return err
}

// SendInterested tells the peer we want pieces from it
func (c *Client) SendInterested() error {
	msg := Message{ID: MsgInterested}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

// SendNotInterested tells the peer we don't want pieces from it
func (c *Client) SendNotInterested() error {
	msg := Message{ID: MsgNotInterested}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

// SendUnchoke tells the peer it can request pieces from us
func (c *Client) SendUnchoke() error {
	msg := Message{ID: MsgUnchoke}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

// SendChoke tells the peer it can't request pieces from us
func (c *Client) SendChoke() error {
	msg := Message{ID: MsgChoke}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

// SendHave tells the peer we have a piece
func (c *Client) SendHave(index int) error {
	msg := Message{ID: MsgHave, Payload: make([]byte, 4)}
	binary.BigEndian.PutUint32(msg.Payload, uint32(index))
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

// recvBitfield reads the bitfield the peer sends after the handshake
func recvBitfield(conn net.Conn) (Bitfield, error) {
	conn.SetDeadline(time.Now().Add(5 * time.Second))
	defer conn.SetDeadline(time.Time{})

	msg, err := ReadMessage(conn)
	if err != nil {
		return nil, err
	}
	if msg == nil {
		return nil, fmt.Errorf("Expected bitfield but got keep-alive")
	}
	if msg.ID != MsgBitfield {
		return nil, fmt.Errorf("Expected bitfield but got ID %d", msg.ID)
	}

	return msg.Payload, nil
}
