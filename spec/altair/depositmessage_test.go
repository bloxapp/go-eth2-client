// Copyright © 2020 Attestant Limited.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package altair_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	spec "github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/goccy/go-yaml"
	require "github.com/stretchr/testify/require"
	"gotest.tools/assert"
)

func TestDepositMessageJSON(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		err   string
	}{
		{
			name: "Empty",
			err:  "unexpected end of JSON input",
		},
		{
			name:  "JSONBad",
			input: []byte("[]"),
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type altair.depositMessageJSON",
		},
		{
			name:  "PublicKeyMissing",
			input: []byte(`{"withdrawal_credentials":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","amount":"32000000000"}`),
			err:   "public key missing",
		},
		{
			name:  "PublicKeyWrongType",
			input: []byte(`{"pubkey":true,"withdrawal_credentials":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","amount":"32000000000"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field depositMessageJSON.pubkey of type string",
		},
		{
			name:  "PublicKeyInvalid",
			input: []byte(`{"pubkey":"invalid","withdrawal_credentials":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","amount":"32000000000"}`),
			err:   "invalid value for public key: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "PublicKeyShort",
			input: []byte(`{"pubkey":"0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","withdrawal_credentials":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","amount":"32000000000"}`),
			err:   "incorrect length for public key",
		},
		{
			name:  "PublicKeyLong",
			input: []byte(`{"pubkey":"0x00000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","withdrawal_credentials":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","amount":"32000000000"}`),
			err:   "incorrect length for public key",
		},
		{
			name:  "WithdrawalCredentialsMissing",
			input: []byte(`{"pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","amount":"32000000000"}`),
			err:   "withdrawal credentials missing",
		},
		{
			name:  "WithdrawalCredentialsWrongType",
			input: []byte(`{"pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","withdrawal_credentials":true,"amount":"32000000000"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field depositMessageJSON.withdrawal_credentials of type string",
		},
		{
			name:  "WithdrawalCredentialsInvalid",
			input: []byte(`{"pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","withdrawal_credentials":"invalid","amount":"32000000000"}`),
			err:   "invalid value for withdrawal credentials: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "WithdrawalCredentialsShort",
			input: []byte(`{"pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","withdrawal_credentials":"0x2122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","amount":"32000000000"}`),
			err:   "incorrect length for withdrawal credentials",
		},
		{
			name:  "WithdrawalCredentialsLong",
			input: []byte(`{"pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","withdrawal_credentials":"0x20202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","amount":"32000000000"}`),
			err:   "incorrect length for withdrawal credentials",
		},
		{
			name:  "AmountMissing",
			input: []byte(`{"pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","withdrawal_credentials":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"}`),
			err:   "amount missing",
		},
		{
			name:  "AmountWrongType",
			input: []byte(`{"pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","withdrawal_credentials":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","amount":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field depositMessageJSON.amount of type string",
		},
		{
			name:  "AmountInvalid",
			input: []byte(`{"pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","withdrawal_credentials":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","amount":"-1"}`),
			err:   "invalid value for amount: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "Good",
			input: []byte(`{"pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","withdrawal_credentials":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","amount":"32000000000"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res spec.DepositMessage
			err := json.Unmarshal(test.input, &res)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				rt, err := json.Marshal(&res)
				require.NoError(t, err)
				assert.Equal(t, string(test.input), string(rt))
			}
		})
	}
}

func TestDepositMessageYAML(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		root  []byte
		err   string
	}{
		{
			name:  "Good",
			input: []byte(`{pubkey: '0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f', withdrawal_credentials: '0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f', amount: 32000000000}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res spec.DepositMessage
			err := yaml.Unmarshal(test.input, &res)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				rt, err := yaml.Marshal(&res)
				require.NoError(t, err)
				rt = bytes.TrimSuffix(rt, []byte("\n"))
				assert.Equal(t, string(test.input), string(rt))
			}
		})
	}
}

func TestDepositMessageSpec(t *testing.T) {
	if os.Getenv("ETH2_SPEC_TESTS_DIR") == "" {
		t.Skip("ETH2_SPEC_TESTS_DIR not suppplied, not running spec tests")
	}
	baseDir := filepath.Join(os.Getenv("ETH2_SPEC_TESTS_DIR"), "tests", "mainnet", "altair", "ssz_static", "DepositMessage", "ssz_random")
	require.NoError(t, filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if path == baseDir {
			// Only interested in subdirectories.
			return nil
		}
		require.NoError(t, err)
		if info.IsDir() {
			t.Run(info.Name(), func(t *testing.T) {
				specYAML, err := ioutil.ReadFile(filepath.Join(path, "value.yaml"))
				require.NoError(t, err)
				var res spec.DepositMessage
				require.NoError(t, yaml.Unmarshal(specYAML, &res))

				specSSZ, err := ioutil.ReadFile(filepath.Join(path, "serialized.ssz"))
				require.NoError(t, err)

				ssz, err := res.MarshalSSZ()
				require.NoError(t, err)
				require.Equal(t, specSSZ, ssz)

				root, err := res.HashTreeRoot()
				require.NoError(t, err)
				rootsYAML, err := ioutil.ReadFile(filepath.Join(path, "roots.yaml"))
				require.NoError(t, err)
				require.Equal(t, string(rootsYAML), fmt.Sprintf("{root: '%#x'}\n", root))
			})
		}
		return nil
	}))
}
