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

package multi

import (
	"context"

	eth2client "github.com/attestantio/go-eth2-client"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
)

// AggregateAttestation fetches the aggregate attestation given an attestation.
func (s *Service) AggregateAttestation(ctx context.Context,
	slot spec.Slot,
	attestationDataRoot spec.Root,
) (
	*spec.Attestation,
	error,
) {
	res, err := s.doCall(ctx, func(ctx context.Context, client eth2client.Service) (interface{}, error) {
		aggregate, err := client.(eth2client.AggregateAttestationProvider).AggregateAttestation(ctx, slot, attestationDataRoot)
		if err != nil {
			return nil, err
		}
		return aggregate, nil
	}, nil)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	return res.(*spec.Attestation), nil
}