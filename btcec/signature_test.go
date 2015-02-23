// Copyright (c) 2013-2014 Conformal Systems LLC.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package btcec_test

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"

	"github.com/melange-app/nmcd/btcec"
)

type signatureTest struct {
	name    string
	sig     []byte
	der     bool
	isValid bool
}

var signatureTests = []signatureTest{
	// signatures from bitcoin blockchain tx
	// 0437cd7f8525ceed2324359c2d0ba26006d92d85
	{
		name: "valid signature.",
		sig: []byte{0x30, 0x44, 0x02, 0x20, 0x4e, 0x45, 0xe1, 0x69,
			0x32, 0xb8, 0xaf, 0x51, 0x49, 0x61, 0xa1, 0xd3, 0xa1,
			0xa2, 0x5f, 0xdf, 0x3f, 0x4f, 0x77, 0x32, 0xe9, 0xd6,
			0x24, 0xc6, 0xc6, 0x15, 0x48, 0xab, 0x5f, 0xb8, 0xcd,
			0x41, 0x02, 0x20, 0x18, 0x15, 0x22, 0xec, 0x8e, 0xca,
			0x07, 0xde, 0x48, 0x60, 0xa4, 0xac, 0xdd, 0x12, 0x90,
			0x9d, 0x83, 0x1c, 0xc5, 0x6c, 0xbb, 0xac, 0x46, 0x22,
			0x08, 0x22, 0x21, 0xa8, 0x76, 0x8d, 0x1d, 0x09,
		},
		der:     true,
		isValid: true,
	},
	{
		name:    "empty.",
		sig:     []byte{},
		isValid: false,
	},
	{
		name: "bad magic.",
		sig: []byte{0x31, 0x44, 0x02, 0x20, 0x4e, 0x45, 0xe1, 0x69,
			0x32, 0xb8, 0xaf, 0x51, 0x49, 0x61, 0xa1, 0xd3, 0xa1,
			0xa2, 0x5f, 0xdf, 0x3f, 0x4f, 0x77, 0x32, 0xe9, 0xd6,
			0x24, 0xc6, 0xc6, 0x15, 0x48, 0xab, 0x5f, 0xb8, 0xcd,
			0x41, 0x02, 0x20, 0x18, 0x15, 0x22, 0xec, 0x8e, 0xca,
			0x07, 0xde, 0x48, 0x60, 0xa4, 0xac, 0xdd, 0x12, 0x90,
			0x9d, 0x83, 0x1c, 0xc5, 0x6c, 0xbb, 0xac, 0x46, 0x22,
			0x08, 0x22, 0x21, 0xa8, 0x76, 0x8d, 0x1d, 0x09,
		},
		der:     true,
		isValid: false,
	},
	{
		name: "bad 1st int marker magic.",
		sig: []byte{0x30, 0x44, 0x03, 0x20, 0x4e, 0x45, 0xe1, 0x69,
			0x32, 0xb8, 0xaf, 0x51, 0x49, 0x61, 0xa1, 0xd3, 0xa1,
			0xa2, 0x5f, 0xdf, 0x3f, 0x4f, 0x77, 0x32, 0xe9, 0xd6,
			0x24, 0xc6, 0xc6, 0x15, 0x48, 0xab, 0x5f, 0xb8, 0xcd,
			0x41, 0x02, 0x20, 0x18, 0x15, 0x22, 0xec, 0x8e, 0xca,
			0x07, 0xde, 0x48, 0x60, 0xa4, 0xac, 0xdd, 0x12, 0x90,
			0x9d, 0x83, 0x1c, 0xc5, 0x6c, 0xbb, 0xac, 0x46, 0x22,
			0x08, 0x22, 0x21, 0xa8, 0x76, 0x8d, 0x1d, 0x09,
		},
		der:     true,
		isValid: false,
	},
	{
		name: "bad 2nd int marker.",
		sig: []byte{0x30, 0x44, 0x02, 0x20, 0x4e, 0x45, 0xe1, 0x69,
			0x32, 0xb8, 0xaf, 0x51, 0x49, 0x61, 0xa1, 0xd3, 0xa1,
			0xa2, 0x5f, 0xdf, 0x3f, 0x4f, 0x77, 0x32, 0xe9, 0xd6,
			0x24, 0xc6, 0xc6, 0x15, 0x48, 0xab, 0x5f, 0xb8, 0xcd,
			0x41, 0x03, 0x20, 0x18, 0x15, 0x22, 0xec, 0x8e, 0xca,
			0x07, 0xde, 0x48, 0x60, 0xa4, 0xac, 0xdd, 0x12, 0x90,
			0x9d, 0x83, 0x1c, 0xc5, 0x6c, 0xbb, 0xac, 0x46, 0x22,
			0x08, 0x22, 0x21, 0xa8, 0x76, 0x8d, 0x1d, 0x09,
		},
		der:     true,
		isValid: false,
	},
	{
		name: "short len",
		sig: []byte{0x30, 0x43, 0x02, 0x20, 0x4e, 0x45, 0xe1, 0x69,
			0x32, 0xb8, 0xaf, 0x51, 0x49, 0x61, 0xa1, 0xd3, 0xa1,
			0xa2, 0x5f, 0xdf, 0x3f, 0x4f, 0x77, 0x32, 0xe9, 0xd6,
			0x24, 0xc6, 0xc6, 0x15, 0x48, 0xab, 0x5f, 0xb8, 0xcd,
			0x41, 0x02, 0x20, 0x18, 0x15, 0x22, 0xec, 0x8e, 0xca,
			0x07, 0xde, 0x48, 0x60, 0xa4, 0xac, 0xdd, 0x12, 0x90,
			0x9d, 0x83, 0x1c, 0xc5, 0x6c, 0xbb, 0xac, 0x46, 0x22,
			0x08, 0x22, 0x21, 0xa8, 0x76, 0x8d, 0x1d, 0x09,
		},
		der:     true,
		isValid: false,
	},
	{
		name: "long len",
		sig: []byte{0x30, 0x45, 0x02, 0x20, 0x4e, 0x45, 0xe1, 0x69,
			0x32, 0xb8, 0xaf, 0x51, 0x49, 0x61, 0xa1, 0xd3, 0xa1,
			0xa2, 0x5f, 0xdf, 0x3f, 0x4f, 0x77, 0x32, 0xe9, 0xd6,
			0x24, 0xc6, 0xc6, 0x15, 0x48, 0xab, 0x5f, 0xb8, 0xcd,
			0x41, 0x02, 0x20, 0x18, 0x15, 0x22, 0xec, 0x8e, 0xca,
			0x07, 0xde, 0x48, 0x60, 0xa4, 0xac, 0xdd, 0x12, 0x90,
			0x9d, 0x83, 0x1c, 0xc5, 0x6c, 0xbb, 0xac, 0x46, 0x22,
			0x08, 0x22, 0x21, 0xa8, 0x76, 0x8d, 0x1d, 0x09,
		},
		der:     true,
		isValid: false,
	},
	{
		name: "long X",
		sig: []byte{0x30, 0x44, 0x02, 0x42, 0x4e, 0x45, 0xe1, 0x69,
			0x32, 0xb8, 0xaf, 0x51, 0x49, 0x61, 0xa1, 0xd3, 0xa1,
			0xa2, 0x5f, 0xdf, 0x3f, 0x4f, 0x77, 0x32, 0xe9, 0xd6,
			0x24, 0xc6, 0xc6, 0x15, 0x48, 0xab, 0x5f, 0xb8, 0xcd,
			0x41, 0x02, 0x20, 0x18, 0x15, 0x22, 0xec, 0x8e, 0xca,
			0x07, 0xde, 0x48, 0x60, 0xa4, 0xac, 0xdd, 0x12, 0x90,
			0x9d, 0x83, 0x1c, 0xc5, 0x6c, 0xbb, 0xac, 0x46, 0x22,
			0x08, 0x22, 0x21, 0xa8, 0x76, 0x8d, 0x1d, 0x09,
		},
		der:     true,
		isValid: false,
	},
	{
		name: "long Y",
		sig: []byte{0x30, 0x44, 0x02, 0x20, 0x4e, 0x45, 0xe1, 0x69,
			0x32, 0xb8, 0xaf, 0x51, 0x49, 0x61, 0xa1, 0xd3, 0xa1,
			0xa2, 0x5f, 0xdf, 0x3f, 0x4f, 0x77, 0x32, 0xe9, 0xd6,
			0x24, 0xc6, 0xc6, 0x15, 0x48, 0xab, 0x5f, 0xb8, 0xcd,
			0x41, 0x02, 0x21, 0x18, 0x15, 0x22, 0xec, 0x8e, 0xca,
			0x07, 0xde, 0x48, 0x60, 0xa4, 0xac, 0xdd, 0x12, 0x90,
			0x9d, 0x83, 0x1c, 0xc5, 0x6c, 0xbb, 0xac, 0x46, 0x22,
			0x08, 0x22, 0x21, 0xa8, 0x76, 0x8d, 0x1d, 0x09,
		},
		der:     true,
		isValid: false,
	},
	{
		name: "short Y",
		sig: []byte{0x30, 0x44, 0x02, 0x20, 0x4e, 0x45, 0xe1, 0x69,
			0x32, 0xb8, 0xaf, 0x51, 0x49, 0x61, 0xa1, 0xd3, 0xa1,
			0xa2, 0x5f, 0xdf, 0x3f, 0x4f, 0x77, 0x32, 0xe9, 0xd6,
			0x24, 0xc6, 0xc6, 0x15, 0x48, 0xab, 0x5f, 0xb8, 0xcd,
			0x41, 0x02, 0x19, 0x18, 0x15, 0x22, 0xec, 0x8e, 0xca,
			0x07, 0xde, 0x48, 0x60, 0xa4, 0xac, 0xdd, 0x12, 0x90,
			0x9d, 0x83, 0x1c, 0xc5, 0x6c, 0xbb, 0xac, 0x46, 0x22,
			0x08, 0x22, 0x21, 0xa8, 0x76, 0x8d, 0x1d, 0x09,
		},
		der:     true,
		isValid: false,
	},
	{
		name: "trailing crap.",
		sig: []byte{0x30, 0x44, 0x02, 0x20, 0x4e, 0x45, 0xe1, 0x69,
			0x32, 0xb8, 0xaf, 0x51, 0x49, 0x61, 0xa1, 0xd3, 0xa1,
			0xa2, 0x5f, 0xdf, 0x3f, 0x4f, 0x77, 0x32, 0xe9, 0xd6,
			0x24, 0xc6, 0xc6, 0x15, 0x48, 0xab, 0x5f, 0xb8, 0xcd,
			0x41, 0x02, 0x20, 0x18, 0x15, 0x22, 0xec, 0x8e, 0xca,
			0x07, 0xde, 0x48, 0x60, 0xa4, 0xac, 0xdd, 0x12, 0x90,
			0x9d, 0x83, 0x1c, 0xc5, 0x6c, 0xbb, 0xac, 0x46, 0x22,
			0x08, 0x22, 0x21, 0xa8, 0x76, 0x8d, 0x1d, 0x09, 0x01,
		},
		der: true,

		// This test is now passing (used to be failing) because there
		// are signatures in the blockchain that have trailing zero
		// bytes before the hashtype. So ParseSignature was fixed to
		// permit buffers with trailing nonsense after the actual
		// signature.
		isValid: true,
	},
	{
		name: "X == N ",
		sig: []byte{0x30, 0x44, 0x02, 0x20, 0xFF, 0xFF, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFE, 0xBA, 0xAE, 0xDC, 0xE6, 0xAF, 0x48,
			0xA0, 0x3B, 0xBF, 0xD2, 0x5E, 0x8C, 0xD0, 0x36, 0x41,
			0x41, 0x02, 0x20, 0x18, 0x15, 0x22, 0xec, 0x8e, 0xca,
			0x07, 0xde, 0x48, 0x60, 0xa4, 0xac, 0xdd, 0x12, 0x90,
			0x9d, 0x83, 0x1c, 0xc5, 0x6c, 0xbb, 0xac, 0x46, 0x22,
			0x08, 0x22, 0x21, 0xa8, 0x76, 0x8d, 0x1d, 0x09,
		},
		der:     true,
		isValid: false,
	},
	{
		name: "X == N ",
		sig: []byte{0x30, 0x44, 0x02, 0x20, 0xFF, 0xFF, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFE, 0xBA, 0xAE, 0xDC, 0xE6, 0xAF, 0x48,
			0xA0, 0x3B, 0xBF, 0xD2, 0x5E, 0x8C, 0xD0, 0x36, 0x41,
			0x42, 0x02, 0x20, 0x18, 0x15, 0x22, 0xec, 0x8e, 0xca,
			0x07, 0xde, 0x48, 0x60, 0xa4, 0xac, 0xdd, 0x12, 0x90,
			0x9d, 0x83, 0x1c, 0xc5, 0x6c, 0xbb, 0xac, 0x46, 0x22,
			0x08, 0x22, 0x21, 0xa8, 0x76, 0x8d, 0x1d, 0x09,
		},
		der:     false,
		isValid: false,
	},
	{
		name: "Y == N",
		sig: []byte{0x30, 0x44, 0x02, 0x20, 0x4e, 0x45, 0xe1, 0x69,
			0x32, 0xb8, 0xaf, 0x51, 0x49, 0x61, 0xa1, 0xd3, 0xa1,
			0xa2, 0x5f, 0xdf, 0x3f, 0x4f, 0x77, 0x32, 0xe9, 0xd6,
			0x24, 0xc6, 0xc6, 0x15, 0x48, 0xab, 0x5f, 0xb8, 0xcd,
			0x41, 0x02, 0x20, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
			0xFE, 0xBA, 0xAE, 0xDC, 0xE6, 0xAF, 0x48, 0xA0, 0x3B,
			0xBF, 0xD2, 0x5E, 0x8C, 0xD0, 0x36, 0x41, 0x41,
		},
		der:     true,
		isValid: false,
	},
	{
		name: "Y > N",
		sig: []byte{0x30, 0x44, 0x02, 0x20, 0x4e, 0x45, 0xe1, 0x69,
			0x32, 0xb8, 0xaf, 0x51, 0x49, 0x61, 0xa1, 0xd3, 0xa1,
			0xa2, 0x5f, 0xdf, 0x3f, 0x4f, 0x77, 0x32, 0xe9, 0xd6,
			0x24, 0xc6, 0xc6, 0x15, 0x48, 0xab, 0x5f, 0xb8, 0xcd,
			0x41, 0x02, 0x20, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
			0xFE, 0xBA, 0xAE, 0xDC, 0xE6, 0xAF, 0x48, 0xA0, 0x3B,
			0xBF, 0xD2, 0x5E, 0x8C, 0xD0, 0x36, 0x41, 0x42,
		},
		der:     false,
		isValid: false,
	},
	{
		name: "0 len X.",
		sig: []byte{0x30, 0x24, 0x02, 0x00, 0x02, 0x20, 0x18, 0x15,
			0x22, 0xec, 0x8e, 0xca, 0x07, 0xde, 0x48, 0x60, 0xa4,
			0xac, 0xdd, 0x12, 0x90, 0x9d, 0x83, 0x1c, 0xc5, 0x6c,
			0xbb, 0xac, 0x46, 0x22, 0x08, 0x22, 0x21, 0xa8, 0x76,
			0x8d, 0x1d, 0x09,
		},
		der:     true,
		isValid: false,
	},
	{
		name: "0 len Y.",
		sig: []byte{0x30, 0x24, 0x02, 0x20, 0x4e, 0x45, 0xe1, 0x69,
			0x32, 0xb8, 0xaf, 0x51, 0x49, 0x61, 0xa1, 0xd3, 0xa1,
			0xa2, 0x5f, 0xdf, 0x3f, 0x4f, 0x77, 0x32, 0xe9, 0xd6,
			0x24, 0xc6, 0xc6, 0x15, 0x48, 0xab, 0x5f, 0xb8, 0xcd,
			0x41, 0x02, 0x00,
		},
		der:     true,
		isValid: false,
	},
	{
		name: "extra R padding.",
		sig: []byte{0x30, 0x45, 0x02, 0x21, 0x00, 0x4e, 0x45, 0xe1, 0x69,
			0x32, 0xb8, 0xaf, 0x51, 0x49, 0x61, 0xa1, 0xd3, 0xa1,
			0xa2, 0x5f, 0xdf, 0x3f, 0x4f, 0x77, 0x32, 0xe9, 0xd6,
			0x24, 0xc6, 0xc6, 0x15, 0x48, 0xab, 0x5f, 0xb8, 0xcd,
			0x41, 0x02, 0x20, 0x18, 0x15, 0x22, 0xec, 0x8e, 0xca,
			0x07, 0xde, 0x48, 0x60, 0xa4, 0xac, 0xdd, 0x12, 0x90,
			0x9d, 0x83, 0x1c, 0xc5, 0x6c, 0xbb, 0xac, 0x46, 0x22,
			0x08, 0x22, 0x21, 0xa8, 0x76, 0x8d, 0x1d, 0x09,
		},
		der:     true,
		isValid: false,
	},
	{
		name: "extra S padding.",
		sig: []byte{0x30, 0x45, 0x02, 0x20, 0x4e, 0x45, 0xe1, 0x69,
			0x32, 0xb8, 0xaf, 0x51, 0x49, 0x61, 0xa1, 0xd3, 0xa1,
			0xa2, 0x5f, 0xdf, 0x3f, 0x4f, 0x77, 0x32, 0xe9, 0xd6,
			0x24, 0xc6, 0xc6, 0x15, 0x48, 0xab, 0x5f, 0xb8, 0xcd,
			0x41, 0x02, 0x21, 0x00, 0x18, 0x15, 0x22, 0xec, 0x8e, 0xca,
			0x07, 0xde, 0x48, 0x60, 0xa4, 0xac, 0xdd, 0x12, 0x90,
			0x9d, 0x83, 0x1c, 0xc5, 0x6c, 0xbb, 0xac, 0x46, 0x22,
			0x08, 0x22, 0x21, 0xa8, 0x76, 0x8d, 0x1d, 0x09,
		},
		der:     true,
		isValid: false,
	},
	// Standard checks (in BER format, without checking for 'canonical' DER
	// signatures) don't test for negative numbers here because there isn't
	// a way that is the same between openssl and go that will mark a number
	// as negative. The Go ASN.1 parser marks numbers as negative when
	// openssl does not (it doesn't handle negative numbers that I can tell
	// at all. When not parsing DER signatures, which is done by by bitcoind
	// when accepting transactions into its mempool, we otherwise only check
	// for the coordinates being zero.
	{
		name: "X == 0",
		sig: []byte{0x30, 0x25, 0x02, 0x01, 0x00, 0x02, 0x20, 0x18,
			0x15, 0x22, 0xec, 0x8e, 0xca, 0x07, 0xde, 0x48, 0x60,
			0xa4, 0xac, 0xdd, 0x12, 0x90, 0x9d, 0x83, 0x1c, 0xc5,
			0x6c, 0xbb, 0xac, 0x46, 0x22, 0x08, 0x22, 0x21, 0xa8,
			0x76, 0x8d, 0x1d, 0x09,
		},
		der:     false,
		isValid: false,
	},
	{
		name: "Y == 0.",
		sig: []byte{0x30, 0x25, 0x02, 0x20, 0x4e, 0x45, 0xe1, 0x69,
			0x32, 0xb8, 0xaf, 0x51, 0x49, 0x61, 0xa1, 0xd3, 0xa1,
			0xa2, 0x5f, 0xdf, 0x3f, 0x4f, 0x77, 0x32, 0xe9, 0xd6,
			0x24, 0xc6, 0xc6, 0x15, 0x48, 0xab, 0x5f, 0xb8, 0xcd,
			0x41, 0x02, 0x01, 0x00,
		},
		der:     false,
		isValid: false,
	},
}

func TestSignatures(t *testing.T) {
	for _, test := range signatureTests {
		var err error
		if test.der {
			_, err = btcec.ParseDERSignature(test.sig, btcec.S256())
		} else {
			_, err = btcec.ParseSignature(test.sig, btcec.S256())
		}
		if err != nil {
			if test.isValid {
				t.Errorf("%s signature failed when shouldn't %v",
					test.name, err)
			} /* else {
				t.Errorf("%s got error %v", test.name, err)
			} */
			continue
		}
		if !test.isValid {
			t.Errorf("%s counted as valid when it should fail",
				test.name)
		}
	}
}

// TestSignatureSerialize ensures that serializing signatures works as expected.
func TestSignatureSerialize(t *testing.T) {
	tests := []struct {
		name     string
		ecsig    *btcec.Signature
		expected []byte
	}{
		// signature from bitcoin blockchain tx
		// 0437cd7f8525ceed2324359c2d0ba26006d92d85
		{
			"valid 1 - r and s most significant bits are zero",
			&btcec.Signature{
				R: fromHex("4e45e16932b8af514961a1d3a1a25fdf3f4f7732e9d624c6c61548ab5fb8cd41"),
				S: fromHex("181522ec8eca07de4860a4acdd12909d831cc56cbbac4622082221a8768d1d09"),
			},
			[]byte{
				0x30, 0x44, 0x02, 0x20, 0x4e, 0x45, 0xe1, 0x69,
				0x32, 0xb8, 0xaf, 0x51, 0x49, 0x61, 0xa1, 0xd3,
				0xa1, 0xa2, 0x5f, 0xdf, 0x3f, 0x4f, 0x77, 0x32,
				0xe9, 0xd6, 0x24, 0xc6, 0xc6, 0x15, 0x48, 0xab,
				0x5f, 0xb8, 0xcd, 0x41, 0x02, 0x20, 0x18, 0x15,
				0x22, 0xec, 0x8e, 0xca, 0x07, 0xde, 0x48, 0x60,
				0xa4, 0xac, 0xdd, 0x12, 0x90, 0x9d, 0x83, 0x1c,
				0xc5, 0x6c, 0xbb, 0xac, 0x46, 0x22, 0x08, 0x22,
				0x21, 0xa8, 0x76, 0x8d, 0x1d, 0x09,
			},
		},
		// signature from bitcoin blockchain tx
		// cb00f8a0573b18faa8c4f467b049f5d202bf1101d9ef2633bc611be70376a4b4
		{
			"valid 2 - r most significant bit is one",
			&btcec.Signature{
				R: fromHex("0082235e21a2300022738dabb8e1bbd9d19cfb1e7ab8c30a23b0afbb8d178abcf3"),
				S: fromHex("24bf68e256c534ddfaf966bf908deb944305596f7bdcc38d69acad7f9c868724"),
			},
			[]byte{
				0x30, 0x45, 0x02, 0x21, 0x00, 0x82, 0x23, 0x5e,
				0x21, 0xa2, 0x30, 0x00, 0x22, 0x73, 0x8d, 0xab,
				0xb8, 0xe1, 0xbb, 0xd9, 0xd1, 0x9c, 0xfb, 0x1e,
				0x7a, 0xb8, 0xc3, 0x0a, 0x23, 0xb0, 0xaf, 0xbb,
				0x8d, 0x17, 0x8a, 0xbc, 0xf3, 0x02, 0x20, 0x24,
				0xbf, 0x68, 0xe2, 0x56, 0xc5, 0x34, 0xdd, 0xfa,
				0xf9, 0x66, 0xbf, 0x90, 0x8d, 0xeb, 0x94, 0x43,
				0x05, 0x59, 0x6f, 0x7b, 0xdc, 0xc3, 0x8d, 0x69,
				0xac, 0xad, 0x7f, 0x9c, 0x86, 0x87, 0x24,
			},
		},
		// signature from bitcoin blockchain tx
		// fda204502a3345e08afd6af27377c052e77f1fefeaeb31bdd45f1e1237ca5470
		{
			"valid 3 - s most significant bit is one",
			&btcec.Signature{
				R: fromHex("1cadddc2838598fee7dc35a12b340c6bde8b389f7bfd19a1252a17c4b5ed2d71"),
				S: new(big.Int).Add(fromHex("00c1a251bbecb14b058a8bd77f65de87e51c47e95904f4c0e9d52eddc21c1415ac"), btcec.S256().N),
			},
			[]byte{
				0x30, 0x45, 0x02, 0x20, 0x1c, 0xad, 0xdd, 0xc2,
				0x83, 0x85, 0x98, 0xfe, 0xe7, 0xdc, 0x35, 0xa1,
				0x2b, 0x34, 0x0c, 0x6b, 0xde, 0x8b, 0x38, 0x9f,
				0x7b, 0xfd, 0x19, 0xa1, 0x25, 0x2a, 0x17, 0xc4,
				0xb5, 0xed, 0x2d, 0x71, 0x02, 0x21, 0x00, 0xc1,
				0xa2, 0x51, 0xbb, 0xec, 0xb1, 0x4b, 0x05, 0x8a,
				0x8b, 0xd7, 0x7f, 0x65, 0xde, 0x87, 0xe5, 0x1c,
				0x47, 0xe9, 0x59, 0x04, 0xf4, 0xc0, 0xe9, 0xd5,
				0x2e, 0xdd, 0xc2, 0x1c, 0x14, 0x15, 0xac,
			},
		},
		{
			"zero signature",
			&btcec.Signature{
				R: big.NewInt(0),
				S: big.NewInt(0),
			},
			[]byte{0x30, 0x06, 0x02, 0x01, 0x00, 0x02, 0x01, 0x00},
		},
	}

	for i, test := range tests {
		result := test.ecsig.Serialize()
		if !bytes.Equal(result, test.expected) {
			t.Errorf("Serialize #%d (%s) unexpected result:\n"+
				"got:  %x\nwant: %x", i, test.name, result,
				test.expected)
		}
	}
}

func testSignCompact(t *testing.T, tag string, curve *btcec.KoblitzCurve,
	data []byte, isCompressed bool) {
	tmp, _ := btcec.NewPrivateKey(curve)
	priv := (*btcec.PrivateKey)(tmp)

	hashed := []byte("testing")
	sig, err := btcec.SignCompact(curve, priv, hashed, isCompressed)
	if err != nil {
		t.Errorf("%s: error signing: %s", tag, err)
		return
	}

	pk, wasCompressed, err := btcec.RecoverCompact(curve, sig, hashed)
	if err != nil {
		t.Errorf("%s: error recovering: %s", tag, err)
		return
	}
	if pk.X.Cmp(priv.X) != 0 || pk.Y.Cmp(priv.Y) != 0 {
		t.Errorf("%s: recovered pubkey doesn't match original "+
			"(%v,%v) vs (%v,%v) ", tag, pk.X, pk.Y, priv.X, priv.Y)
		return
	}
	if wasCompressed != isCompressed {
		t.Errorf("%s: recovered pubkey doesn't match compressed state "+
			"(%v vs %v)", tag, isCompressed, wasCompressed)
		return
	}

	// If we change the compressed bit we should get the same key back,
	// but the compressed flag should be reversed.
	if isCompressed {
		sig[0] -= 4
	} else {
		sig[0] += 4
	}

	pk, wasCompressed, err = btcec.RecoverCompact(curve, sig, hashed)
	if err != nil {
		t.Errorf("%s: error recovering (2): %s", tag, err)
		return
	}
	if pk.X.Cmp(priv.X) != 0 || pk.Y.Cmp(priv.Y) != 0 {
		t.Errorf("%s: recovered pubkey (2) doesn't match original "+
			"(%v,%v) vs (%v,%v) ", tag, pk.X, pk.Y, priv.X, priv.Y)
		return
	}
	if wasCompressed == isCompressed {
		t.Errorf("%s: recovered pubkey doesn't match reversed "+
			"compressed state (%v vs %v)", tag, isCompressed,
			wasCompressed)
		return
	}
}

func TestSignCompact(t *testing.T) {
	for i := 0; i < 256; i++ {
		name := fmt.Sprintf("test %d", i)
		data := make([]byte, 32)
		_, err := rand.Read(data)
		if err != nil {
			t.Errorf("failed to read random data for %s", name)
			continue
		}
		compressed := i%2 != 0
		testSignCompact(t, name, btcec.S256(), data, compressed)
	}
}
