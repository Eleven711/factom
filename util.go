// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"bytes"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/binary"
	"fmt"
	"time"
)

const (
	ZeroHash = "0000000000000000000000000000000000000000000000000000000000000000"
)

var (
	factomdServer = "localhost:8088"
	walletServer  = "localhost:8089"
	RpcConfig     = &RPCConfig{}
)

// SetFactomdServer sets the gloabal target for the factomd server
func SetFactomdServer(s string) {
	factomdServer = s
}

// SetWalletServer sets the global target for the fctwallet server
func SetWalletServer(s string) {
	walletServer = s
}

// FactomdServer returns the global server string for debugging
func FactomdServer() string {
	return factomdServer
}

// FactomdServer returns the global wallet server string for debugging
func WalletServer() string {
	return walletServer
}

// milliTime returns a 6 byte slice representing the unix time in milliseconds
func milliTime() (r []byte) {
	buf := new(bytes.Buffer)
	t := time.Now().UnixNano()
	m := t / 1e6
	binary.Write(buf, binary.BigEndian, m)
	return buf.Bytes()[2:]
}

// shad Double Sha256 Hash; sha256(sha256(data))
func shad(data []byte) []byte {
	h1 := sha256.Sum256(data)
	h2 := sha256.Sum256(h1[:])
	return h2[:]
}

// sha52 Sha512+Sha256 Hash; sha256(sha512(data)+data)
func sha52(data []byte) []byte {
	h1 := sha512.Sum512(data)
	h2 := sha256.Sum256(append(h1[:], data...))
	return h2[:]
}

func entryCost(e *Entry) (int8, error) {
	p, err := e.MarshalBinary()
	if err != nil {
		return 0, err
	}

	// caulculate the length exluding the header size 35 for Milestone 1
	l := len(p) - 35

	if l > 10240 {
		return 10, fmt.Errorf("Entry cannot be larger than 10KB")
	}

	// n is the capacity of the entry payment in KB
	n := int8(l / 1024)

	if r := l % 1024; r > 0 {
		n += 1
	}

	if n < 1 {
		n = 1
	}
	return n, nil
}
