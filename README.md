# Quentin Torrentino

> A (very) simple BitTorrent client built using GO

## Idea
- Parse the `.torrent` file and extract the metadata & `info_hash`.
- Send an "announce" request with the `info_hash` to the tracker to recieve list of peers.
- Establish a TCP connection and perform the BitTorrent handshake with peers.
- Request file "pieces" from peers (seeding not supported yet).
- Verify each received piece using its SHA-1 hash.
- Combine all verified pieces into the final file.

tl;dr: parse -> connect -> download -> verify -> assemble

## `.torrent` file structure
> Contains the torrent metadata and the information to connect to the trackers, encoded in **Bencode**

### Bencode format

- Strings
  - \<string_len>:\<string>
  - eg. 7:comment
- Integers
  - i\<number>e
  - eg. i2e
- Array
  - l\<items>e
  - eg. l5:helloi3ee -> ["hello", 3]
- Dict
  - d\<items>e
  - eg. d3:byei57ee -> {"bye": 57}

The important contents of the `.torrent` file are the announce url (tracker), and a big binary blob which holds the SHA-1 hashes of all pieces. 

This project uses a [pre-built parser](https://github.com/jackpal/bencode-go).

## BitTorrent specification

https://www.bittorrent.org/beps/bep_0003.html

## Inspiration
[Jesse Li](https://blog.jse.li/)