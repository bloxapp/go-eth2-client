// Copyright © 2021 Attestant Limited.
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

package altair

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

// Altair constants.
var syncCommitteeSize = 512

// SyncCommittee is the Ethereum 2 sync committee structure.
type SyncCommittee struct {
	PubKeys         []BLSPubKey `ssz-size:"512,48"`
	AggregatePubKey BLSPubKey
}

// syncCommitteeJSON is the spec representation of the struct.
type syncCommitteeJSON struct {
	PubKeys         []string `json:"pubkeys"`
	AggregatePubKey string   `json:"aggregate_pubkey"`
}

// syncCommitteeYAML is the spec representation of the struct.
type syncCommitteeYAML struct {
	PubKeys         []string `yaml:"pubkeys"`
	AggregatePubKey string   `yaml:"aggregate_pubkey"`
}

// MarshalJSON implements json.Marshaler.
func (s *SyncCommittee) MarshalJSON() ([]byte, error) {
	pubKeys := make([]string, len(s.PubKeys))
	for i := range s.PubKeys {
		pubKeys[i] = fmt.Sprintf("%#x", s.PubKeys[i])
	}

	return json.Marshal(&syncCommitteeJSON{
		PubKeys:         pubKeys,
		AggregatePubKey: fmt.Sprintf("%#x", s.AggregatePubKey),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *SyncCommittee) UnmarshalJSON(input []byte) error {
	var syncCommitteeJSON syncCommitteeJSON
	if err := json.Unmarshal(input, &syncCommitteeJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	return s.unpack(&syncCommitteeJSON)
}

func (s *SyncCommittee) unpack(syncCommitteeJSON *syncCommitteeJSON) error {
	if len(syncCommitteeJSON.PubKeys) != syncCommitteeSize {
		return errors.New("incorrect length for public keys")
	}
	s.PubKeys = make([]BLSPubKey, len(syncCommitteeJSON.PubKeys))
	for i := range syncCommitteeJSON.PubKeys {
		pubKey, err := hex.DecodeString(strings.TrimPrefix(syncCommitteeJSON.PubKeys[i], "0x"))
		if err != nil {
			return errors.Wrap(err, "invalid value for public key")
		}
		if len(pubKey) != PublicKeyLength {
			return errors.New("incorrect length for public key")
		}
		copy(s.PubKeys[i][:], pubKey)
	}

	if syncCommitteeJSON.AggregatePubKey == "" {
		return errors.New("aggregate public key missing")
	}
	aggregatePubKey, err := hex.DecodeString(strings.TrimPrefix(syncCommitteeJSON.AggregatePubKey, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for aggregate public key")
	}
	if len(aggregatePubKey) < 48 {
		return errors.New("aggregate public key short")
	}
	if len(aggregatePubKey) > 48 {
		return errors.New("aggregate public key long")
	}
	copy(s.AggregatePubKey[:], aggregatePubKey)

	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (s *SyncCommittee) MarshalYAML() ([]byte, error) {
	pubKeys := make([]string, len(s.PubKeys))
	for i := range s.PubKeys {
		pubKeys[i] = fmt.Sprintf("%#x", s.PubKeys[i])
	}

	yamlBytes, err := yaml.MarshalWithOptions(&syncCommitteeYAML{
		PubKeys:         pubKeys,
		AggregatePubKey: fmt.Sprintf("%#x", s.AggregatePubKey),
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}
	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (s *SyncCommittee) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var syncCommitteeJSON syncCommitteeJSON
	if err := yaml.Unmarshal(input, &syncCommitteeJSON); err != nil {
		return err
	}
	return s.unpack(&syncCommitteeJSON)
}

// String returns a string version of the structure.
func (s *SyncCommittee) String() string {
	data, err := yaml.Marshal(s)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}
	return string(data)
}
