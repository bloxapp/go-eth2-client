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

package spec

import (
	"errors"

	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// VersionedSignedBeaconBlock contains a versioned signed beacon block.
type VersionedSignedBeaconBlock struct {
	Version DataVersion
	Phase0  *phase0.SignedBeaconBlock
	Altair  *altair.SignedBeaconBlock
}

// Slot returns the slot of the signed beacon block.
func (v *VersionedSignedBeaconBlock) Slot() (phase0.Slot, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil || v.Phase0.Message == nil {
			return 0, errors.New("no phase0 block")
		}
		return v.Phase0.Message.Slot, nil
	case DataVersionAltair:
		if v.Altair == nil || v.Altair.Message == nil {
			return 0, errors.New("no altair block")
		}
		return v.Altair.Message.Slot, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// Attestations returns the attestations of the beacon block.
func (v *VersionedSignedBeaconBlock) Attestations() ([]*phase0.Attestation, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil || v.Phase0.Message == nil || v.Phase0.Message.Body == nil {
			return nil, errors.New("no phase0 block")
		}
		return v.Phase0.Message.Body.Attestations, nil
	case DataVersionAltair:
		if v.Altair == nil || v.Altair.Message == nil || v.Altair.Message.Body == nil {
			return nil, errors.New("no altair block")
		}
		return v.Altair.Message.Body.Attestations, nil
	default:
		return nil, errors.New("unknown version")
	}
}
