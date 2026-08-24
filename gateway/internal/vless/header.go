package vless

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"net"
)

type Header struct {
	UUID    [16]byte
	Command byte
	Port    uint16
	Address string
}

func ParseHeader(r io.Reader) (*Header, error) {
	var head [18]byte // version + uuid + addonLen
	if _, err := io.ReadFull(r, head[:]); err != nil {
		return nil, err
	}
	if head[0] != 0 {
		return nil, fmt.Errorf("unsupported version: %d", head[0])
	}
	if n := int64(head[17]); n > 0 {
		if _, err := io.CopyN(io.Discard, r, n); err != nil {
			return nil, err
		}
	}

	var meta [4]byte // cmd + port + atype
	if _, err := io.ReadFull(r, meta[:]); err != nil {
		return nil, err
	}

	addr, err := readAddress(r, meta[3])
	if err != nil {
		return nil, err
	}

	h := Header{
		Command: meta[0],
		Port:    binary.BigEndian.Uint16(meta[1:3]),
		Address: addr,
	}
	copy(h.UUID[:], head[1:17])
	return &h, nil
}

func readAddress(r io.Reader, atype byte) (string, error) {
	switch atype {
	case 1: // IPv4
		var ip [4]byte
		if _, err := io.ReadFull(r, ip[:]); err != nil {
			return "", err
		}
		return net.IP(ip[:]).String(), nil
	case 2: // Domain
		var l [1]byte
		if _, err := io.ReadFull(r, l[:]); err != nil {
			return "", err
		}
		d := make([]byte, l[0]) // в 1.26 чаще уйдёт на стек
		if _, err := io.ReadFull(r, d); err != nil {
			return "", err
		}
		return string(d), nil
	case 3: // IPv6
		var ip [16]byte
		if _, err := io.ReadFull(r, ip[:]); err != nil {
			return "", err
		}
		return net.IP(ip[:]).String(), nil
	default:
		return "", fmt.Errorf("invalid address type: %d", atype)
	}
}

func FormatUUID(id [16]byte) string {
	b := id[:]
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		hex.EncodeToString(b[0:4]),
		hex.EncodeToString(b[4:6]),
		hex.EncodeToString(b[6:8]),
		hex.EncodeToString(b[8:10]),
		hex.EncodeToString(b[10:16]),
	)
}
