// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package wallet

import (
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"

	"github.com/FactomProject/btcutil/base58"
	"github.com/FactomProject/factom"
	"github.com/FactomProject/factomd/common/factoid"
	"github.com/FactomProject/factomd/common/primitives"
	"github.com/FactomProject/goleveldb/leveldb"
)

var (
	ErrFeeTooHigh    = errors.New("wallet: Overpaying Fee")
	ErrFeeTooLow     = errors.New("wallet: Insufficient Fee")
	ErrNoSuchAddress = errors.New("wallet: No such address")
	ErrTXExists      = errors.New("wallet: Transaction name already exists")
	ErrTXNotExists   = errors.New("wallet: Transaction name was not found")
	ErrTXNoInputs    = errors.New("wallet: Transaction has no inputs")
	ErrTXInvalidName = errors.New("wallet: Transaction name is not valid")
)

func (w *Wallet) NewTransaction(name string) error {
	if _, exist := w.transactions[name]; exist {
		return ErrTXExists
	}

	// check that the transaction name is valid
	if name == "" {
		return ErrTXInvalidName
	}
	if len(name) > 32 {
		return ErrTXInvalidName
	}
	if match, err := regexp.MatchString("[^a-zA-Z0-9_-]", name); err != nil {
		return err
	} else if match {
		return ErrTXInvalidName
	}

	t := new(factoid.Transaction)
	t.SetTimestamp(primitives.NewTimestampNow())
	w.transactions[name] = t
	return nil
}

func (w *Wallet) DeleteTransaction(name string) error {
	if _, exists := w.transactions[name]; !exists {
		return ErrTXNotExists
	}
	delete(w.transactions, name)
	return nil
}

func (w *Wallet) AddInput(name, address string, amount uint64) error {
	if _, exists := w.transactions[name]; !exists {
		return ErrTXNotExists
	}
	trans := w.transactions[name]

	a, err := w.GetFCTAddress(address)
	if err == leveldb.ErrNotFound {
		return ErrNoSuchAddress
	} else if err != nil {
		return err
	}
	adr := factoid.NewAddress(a.RCDHash())

	// First look if this is really an update
	for _, input := range trans.GetInputs() {
		if input.GetAddress().IsSameAs(adr) {
			input.SetAmount(amount)
			return nil
		}
	}

	// Add our new input
	trans.AddInput(adr, amount)
	trans.AddRCD(factoid.NewRCD_1(a.PubBytes()))

	return nil
}

func (w *Wallet) AddOutput(name, address string, amount uint64) error {
	if _, exists := w.transactions[name]; !exists {
		return ErrTXNotExists
	}
	trans := w.transactions[name]

	if !factom.IsValidAddress(address) {
		return errors.New("Invalid Address")
	}

	adr := factoid.NewAddress(base58.Decode(address)[2:34])

	// First look if this is really an update
	for _, output := range trans.GetOutputs() {
		if output.GetAddress().IsSameAs(adr) {
			output.SetAmount(amount)
			return nil
		}
	}

	trans.AddOutput(adr, amount)

	return nil
}

func (w *Wallet) AddECOutput(name, address string, amount uint64) error {
	if _, exists := w.transactions[name]; !exists {
		return ErrTXNotExists
	}
	trans := w.transactions[name]

	if !factom.IsValidAddress(address) {
		return errors.New("Invalid Address")
	}

	adr := factoid.NewAddress(base58.Decode(address)[2:34])

	// First look if this is really an update
	for _, output := range trans.GetECOutputs() {
		if output.GetAddress().IsSameAs(adr) {
			output.SetAmount(amount)
			return nil
		}
	}

	trans.AddECOutput(adr, amount)

	return nil
}

func (w *Wallet) AddFee(name, address string, rate uint64) error {
	if _, exists := w.transactions[name]; !exists {
		return ErrTXNotExists
	}
	trans := w.transactions[name]

	{
		ins, err := trans.TotalInputs()
		if err != nil {
			return err
		}
		outs, err := trans.TotalOutputs()
		if err != nil {
			return err
		}
		ecs, err := trans.TotalECs()
		if err != nil {
			return err
		}

		if ins != outs+ecs {
			return fmt.Errorf("Inputs and outputs don't add up")
		}
	}

	transfee, err := trans.CalculateFee(rate)
	if err != nil {
		return err
	}

	a, err := w.GetFCTAddress(address)
	if err != nil {
		return err
	}
	adr := factoid.NewAddress(a.RCDHash())

	for _, input := range trans.GetInputs() {
		if input.GetAddress().IsSameAs(adr) {
			amt, err := factoid.ValidateAmounts(input.GetAmount(), transfee)
			if err != nil {
				return err
			}
			input.SetAmount(amt)
			return nil
		}
	}
	return fmt.Errorf("%s is not an input to the transaction.", address)
}

func (w *Wallet) SubFee(name, address string, rate uint64) error {
	if _, exists := w.transactions[name]; !exists {
		return ErrTXNotExists
	}
	trans := w.transactions[name]

	if !factom.IsValidAddress(address) {
		return errors.New("Invalid Address")
	}

	{
		ins, err := trans.TotalInputs()
		if err != nil {
			return err
		}
		outs, err := trans.TotalOutputs()
		if err != nil {
			return err
		}
		ecs, err := trans.TotalECs()
		if err != nil {
			return err
		}

		if ins != outs+ecs {
			return fmt.Errorf("Inputs and outputs don't add up")
		}
	}

	transfee, err := trans.CalculateFee(rate)
	if err != nil {
		return err
	}

	adr := factoid.NewAddress(base58.Decode(address)[2:34])

	for _, output := range trans.GetOutputs() {
		if output.GetAddress().IsSameAs(adr) {
			output.SetAmount(output.GetAmount() - transfee)
			return nil
		}
	}
	return fmt.Errorf("%s is not an output to the transaction.", address)
}

func (w *Wallet) SignTransaction(name string) error {
	if _, exists := w.transactions[name]; !exists {
		return ErrTXNotExists
	}
	trans := w.transactions[name]

	if err := checkFee(trans); err != nil {
		return err
	}

	data, err := trans.MarshalBinarySig()
	if err != nil {
		return err
	}

	rcds := trans.GetRCDs()
	if len(rcds) == 0 {
		return ErrTXNoInputs
	}
	for i, rcd := range rcds {
		a, err := rcd.GetAddress()
		if err != nil {
			return err
		}

		f, err := w.GetFCTAddress(primitives.ConvertFctAddressToUserStr(a))
		if err != nil {
			return err
		}
		sig := factoid.NewSingleSignatureBlock(f.SecBytes(), data)
		trans.SetSignatureBlock(i, sig)
	}

	return nil
}

func (w *Wallet) GetTransactions() map[string]*factoid.Transaction {
	return w.transactions
}

func (w *Wallet) ComposeTransaction(name string) (*factom.JSON2Request, error) {
	if _, exists := w.transactions[name]; !exists {
		return nil, ErrTXNotExists
	}
	trans := w.transactions[name]

	type txreq struct {
		Transaction string `json:"transaction"`
	}

	param := new(txreq)
	if p, err := trans.MarshalBinary(); err != nil {
		return nil, err
	} else {
		param.Transaction = hex.EncodeToString(p)
	}

	req := factom.NewJSON2Request("factoid-submit", APICounter(), param)

	return req, nil
}

func checkFee(t *factoid.Transaction) error {
	ins, err := t.TotalInputs()
	if err != nil {
		return err
	}
	outs, err := t.TotalOutputs()
	if err != nil {
		return err
	}
	ecs, err := t.TotalECs()
	if err != nil {
		return err
	}

	// fee is the fee that will be paid
	fee := int64(ins) - int64(outs) - int64(ecs)

	if fee <= 0 {
		return ErrFeeTooLow
	}

	rate, err := factom.GetRate()
	if err != nil {
		return err
	}

	// cfee is the fee calculated for the transaction
	var cfee int64
	if c, err := t.CalculateFee(rate); err != nil {
		return err
	} else if c == 0 {
		return errors.New("wallet: Could not calculate fee")
	} else {
		cfee = int64(c)
	}

	// fee is too low
	if fee < cfee {
		return ErrFeeTooLow
	}

	// fee is too high (over 10x cfee)
	if fee >= cfee*10 {
		return ErrFeeTooHigh
	}

	return nil
}
