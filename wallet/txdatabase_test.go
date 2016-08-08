// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package wallet_test

import (
	. "github.com/FactomProject/factom/wallet"
	
	"os"
	"testing"
	
	"github.com/FactomProject/factom"
	"github.com/FactomProject/factomd/common/factoid"
	"github.com/FactomProject/factomd/common/interfaces"
	
	"fmt" // DEBUG
)

func TestTXDatabaseOverlay(t *testing.T) {
	dbpath := os.TempDir() + "/test_txdb-01"
	db1, err := NewTXLevelDB(dbpath)
	if err != nil {
		t.Error(err)
	}
	defer db1.Close()
	
	fblock, err := fblockHead()
	if err != nil {
		t.Error(err)
	}	
	if err := db1.InsertFBlock(fblock); err != nil {
		t.Error(err)
	}	
	if f, err := db1.GetFBlock(fblock.GetKeyMR().String()); err != nil {
		t.Error(err)
	} else if f == nil {
		t.Errorf("Fblock not found in db")
	}
}

func TestGetAllTXs(t *testing.T) {
	dbpath := os.TempDir() + "/test_txdb-01"
	db1, err := NewTXLevelDB(dbpath)
	if err != nil {
		t.Error(err)
	}
	defer db1.Close()
	
	txs := make(chan interfaces.ITransaction, 500)
	errs := make(chan error)
	
	fmt.Println("DEBUG: running getalltxs")
	go db1.GetAllTXs(txs, errs)
	
	for {
		fmt.Println("DEBUG: started select loop")
		select {
		case tx, ok := <-txs:
			fmt.Println("Got TX:", tx)
			if !ok {
				txs = nil
			}
		case err, ok := <-errs:
			fmt.Println("DEBUG: got error:", err)
			t.Error(err)
			if !ok {
				errs = nil
			}
		}
		
		if txs == nil && errs == nil {
			break
		}
	}
}

// fblockHead gets the most recent fblock.
func fblockHead() (interfaces.IFBlock, error) {
		fblockID := "000000000000000000000000000000000000000000000000000000000000000f"

		dbhead, err := factom.GetDBlockHead()
		if err != nil {
			return nil, err
		}
		dblock, err := factom.GetDBlock(dbhead)
		if err != nil {
			return nil, err
		}
		
		var fblockmr string
		for _, eblock := range dblock.EntryBlockList {
			if eblock.ChainID == fblockID {
				fblockmr = eblock.KeyMR
			}
		}
		if fblockmr == "" {
			return nil, err
		}
		
		// get the most recent block
		p, err := factom.GetRaw(fblockmr)
		if err != nil {
			return nil, err
		}
		return factoid.UnmarshalFBlock(p)
}