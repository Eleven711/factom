// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom_test

import (
	ed "github.com/FactomProject/ed25519"
	"testing"

	. "github.com/FactomProject/factom"
)

var ()

func TestAddressStringType(t *testing.T) {
	var (
		a0 = "FX1zT4aFpEvcnPqPCigB3fvGu4Q4mTXY22iiuV69DqE1pNhdF2MX"
		a1 = "FA1zT4aFpEvcnPqPCigB3fvGu4Q4mTXY22iiuV69DqE1pNhdF2MC"
		a2 = "Fs1KWJrpLdfucvmYwN2nWrwepLn8ercpMbzXshd1g8zyhKXLVLWj"
		a3 = "EC2DKSYyRcNWf7RS963VFYgMExoHRYLHVeCfQ9PGPmNzwrcmgm2r"
		a4 = "Es2Rf7iM6PdsqfYCo3D1tnAR65SkLENyWJG1deUzpRMQmbh9F3eG"
	)

	if v := AddressStringType(a0); v != InvalidAddress {
		t.Errorf("invalid address has wrong type %s %#v", a0, v)
	}
	if v := AddressStringType(a1); v != FactoidPub {
		t.Errorf("wrong address type %s %#v", a1, v)
	}
	if v := AddressStringType(a2); v != FactoidSec {
		t.Errorf("wrong address type %s %#v", a1, v)
	}
	if v := AddressStringType(a3); v != ECPub {
		t.Errorf("wrong address type %s %#v", a1, v)
	}
	if v := AddressStringType(a4); v != ECSec {
		t.Errorf("wrong address type %s %#v", a1, v)
	}
}

func TestNewECAddress(t *testing.T) {
	zPub := "EC1m9mouvUQeEidmqpUYpYtXg8fvTYi6GNHaKg8KMLbdMBrFfmUa"
	e := NewECAddress()
	if e.PubString() != zPub {
		t.Errorf("new address %s did not match %s", e.PubString(), zPub)
	}
}

func TestECAddress(t *testing.T) {
	zPub := "EC1m9mouvUQeEidmqpUYpYtXg8fvTYi6GNHaKg8KMLbdMBrFfmUa"
	zSec := "Es2Rf7iM6PdsqfYCo3D1tnAR65SkLENyWJG1deUzpRMQmbh9F3eG"
	e := NewECAddress()
	e.Pub = &[32]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	e.Sec = &[64]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01}

	if e.PubString() != zPub {
		t.Errorf("%s did not match %s", e.PubString(), zPub)
	}

	if e.SecString() != zSec {
		t.Errorf("%s did not match %s", e.SecString(), zSec)
	}
}

func TestIsValidECAddress(t *testing.T) {
	zPub := "EC1m9mouvUQeEidmqpUYpYtXg8fvTYi6GNHaKg8KMLbdMBrFfmUa"
	zSec := "Es2Rf7iM6PdsqfYCo3D1tnAR65SkLENyWJG1deUzpRMQmbh9F3eG"
	badEmpty := ""
	badLen := "EC1m9mouvUQeEidmqpUYpYtXgfvTYi6GNHaKg8KMLbdMBrFfmUa"
	badPrePub := "Ec1m9mouvUQeEidmqpUYpYtXg8fvTYi6GNHaKg8KMLbdMBrFfmUa"
	badPreSec := "ER2Rf7iM6PdsqfYCo3D1tnAR65SkLENyWJG1deUzpRMQmbh9F3eG"
	badCheckPub := "EC1m9mouvUQeEidmqpUYpYtXg8fvTYi6GNHaKg8KMLbdMBrFfgUa"
	badCheckSec := "Es2Rf7iM6PdsqfYCo3D1tnAR65SkLENyWJG1deUzpRMQmbh9F3ea"

	if !IsValidAddress(zPub) {
		t.Errorf("%s was not considered valid", zPub)
	}
	if !IsValidAddress(zSec) {
		t.Errorf("%s was not considered valid", zSec)
	}

	if IsValidAddress(badEmpty) {
		t.Errorf("%s was considered valid", badEmpty)
	}
	if IsValidAddress(badLen) {
		t.Errorf("%s was considered valid", badLen)
	}
	if IsValidAddress(badPrePub) {
		t.Errorf("%s was considered valid", badPrePub)
	}
	if IsValidAddress(badPreSec) {
		t.Errorf("%s was considered valid", badPreSec)
	}
	if IsValidAddress(badCheckPub) {
		t.Errorf("%s was considered valid", badCheckPub)
	}
	if IsValidAddress(badCheckSec) {
		t.Errorf("%s was considered valid", badCheckSec)
	}
}

func TestGetECAddress(t *testing.T) {
	zSec := "Es2Rf7iM6PdsqfYCo3D1tnAR65SkLENyWJG1deUzpRMQmbh9F3eG"
	e, err := GetECAddress(zSec)
	if err != nil {
		t.Error(err)
	}

	// verify that the keys work
	msg := []byte("Hello Factom!")
	sig := ed.Sign(e.SecFixed(), msg)
	if !ed.Verify(e.PubFixed(), msg, sig) {
		t.Errorf("Key signature did not match")
	}
}

func TestIsValidFactoidAddress(t *testing.T) {
	zPub := "EC1m9mouvUQeEidmqpUYpYtXg8fvTYi6GNHaKg8KMLbdMBrFfmUa"
	zSec := "Es2Rf7iM6PdsqfYCo3D1tnAR65SkLENyWJG1deUzpRMQmbh9F3eG"
	badEmpty := ""
	badLen := "FA1y5ZGuHSLmf2TqNf6hVMkPiNGyQpQDFJvDLRkKQaoPo4bmbgu"
	badPrePub := "Fe1y5ZGuHSLmf2TqNf6hVMkPiNGyQpQDTFJvDLRkKQaoPo4bmbgu"
	badPreSec := "Fb1KWJrpLdfucvmYwN2nWrwepLn8ercpMbzXshd1g8zyhKXLVLWj"
	badCheckPub := "FA1y5ZGuHSLmf2TqNf6hVMkPiNGyQpQDTFJvDLRkKQaoPo4bmggu"
	badCheckSec := "Fs1KWJrpLdfucvmYwN2nWrwepLn8ercpMbzXshd1g8zyhKXLVLwj"

	if !IsValidAddress(zPub) {
		t.Errorf("%s was not considered valid", zPub)
	}
	if !IsValidAddress(zSec) {
		t.Errorf("%s was not considered valid", zSec)
	}

	if IsValidAddress(badEmpty) {
		t.Errorf("%s was considered valid", badEmpty)
	}
	if IsValidAddress(badLen) {
		t.Errorf("%s was considered valid", badLen)
	}
	if IsValidAddress(badPrePub) {
		t.Errorf("%s was considered valid", badPrePub)
	}
	if IsValidAddress(badPreSec) {
		t.Errorf("%s was considered valid", badPreSec)
	}
	if IsValidAddress(badCheckPub) {
		t.Errorf("%s was considered valid", badCheckPub)
	}
	if IsValidAddress(badCheckSec) {
		t.Errorf("%s was considered valid", badCheckSec)
	}
}

func TestGetFactoidAddress(t *testing.T) {
	zSec := "Fs1KWJrpLdfucvmYwN2nWrwepLn8ercpMbzXshd1g8zyhKXLVLWj"

	if _, err := GetFactoidAddress(zSec); err != nil {
		t.Error(err)
	}

	// ? test factoid key validity here
}

func TestMakeFactoidAddressFromMnemonic(t *testing.T) {
	m := "yellow yellow yellow yellow yellow yellow yellow yellow yellow yellow yellow yellow"
	cannonAdr := "FA3cih2o2tjEUsnnFR4jX1tQXPpSXFwsp3rhVp6odL5PNCHWvZV1"

	fct, err := MakeFactoidAddressFromMnemonic(m)
	if err != nil {
		t.Error(err)
	}

	if fct.String() != cannonAdr {
		t.Errorf(
			"incorrect factoid address from 12 words: got %s expecting %s",
			fct.String(), cannonAdr)
	}
}
