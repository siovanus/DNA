/*
 * Copyright (C) 2018 The DNA Authors
 * This file is part of The DNA library.
 *
 * The DNA is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The DNA is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The DNA.  If not, see <http://www.gnu.org/licenses/>.
 */
package ontid

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/dnaproject2/DNA/common/serialization"
	"github.com/dnaproject2/DNA/smartcontract/service/native"
	"github.com/dnaproject2/DNA/smartcontract/service/native/utils"
)

// Group defines a group control logic
type Group struct {
	Members   []interface{} `json:"members"`
	Threshold uint          `json:"threshold"`
}

func (g *Group) ToJson() []byte {
	j, _ := json.Marshal(g)
	return j
}

func deserializeGroup(data []byte) (*Group, error) {
	g := Group{}
	buf := bytes.NewBuffer(data)

	// parse members
	num, err := utils.ReadVarUint(buf)
	if err != nil {
		return nil, fmt.Errorf("error parsing group members: %s", err)
	}

	for i := uint64(0); i < num; i++ {
		m, err := serialization.ReadVarBytes(buf)
		if err != nil {
			return nil, fmt.Errorf("error parsing group members: %s", err)
		}
		if bytes.Equal(m[:8], []byte("did:dna:")) {
			g.Members = append(g.Members, string(m))
		} else {
			// parse recursively
			g1, err := deserializeGroup(m)
			if err != nil {
				return nil, fmt.Errorf("error parsing group members: %s", err)
			}
			g.Members[i] = g1
		}
	}

	// parse threshold
	t, err := utils.ReadVarUint(buf)
	if err != nil {
		return nil, fmt.Errorf("error parsing group threshold: %s", err)
	}

	g.Threshold = uint(t)

	return &g, nil
}

func validateMembers(srvc *native.NativeService, g *Group) error {
	for _, m := range g.Members {
		switch t := m.(type) {
		case []byte:
			key, err := encodeID(t)
			if err != nil {
				return fmt.Errorf("invalid id: %s", string(t))
			}
			// ID must exists
			if !checkIDExistence(srvc, key) {
				return fmt.Errorf("id %s not registered", string(t))
			}
			// Group member must have its own public key
			pk, err := getPk(srvc, key, 1)
			if err != nil || pk == nil {
				return fmt.Errorf("id %s has no public keys", string(t))
			}
		case *Group:
			if err := validateMembers(srvc, t); err != nil {
				return err
			}
		default:
			panic("group member type error")
		}
	}
	return nil
}

type Signer struct {
	id    []byte
	index uint32
}

func deserializeSigners(data []byte) ([]Signer, error) {
	buf := bytes.NewBuffer(data)
	num, err := utils.ReadVarUint(buf)
	if err != nil {
		return nil, err
	}

	signers := []Signer{}
	for i := uint64(0); i < num; i++ {
		id, err := serialization.ReadVarBytes(buf)
		if err != nil {
			return nil, err
		}
		index, err := utils.ReadVarUint(buf)
		if err != nil {
			return nil, err
		}

		signer := Signer{id, uint32(index)}
		signers = append(signers, signer)
	}

	return signers, nil
}

func findSigner(id []byte, signers []Signer) bool {
	for _, signer := range signers {
		if bytes.Equal(signer.id, id) {
			return true
		}
	}
	return false
}

func verifyThreshold(g *Group, signers []Signer) bool {
	var signed uint = 0
	for _, member := range g.Members {
		switch t := member.(type) {
		case []byte:
			if findSigner(t, signers) {
				signed += 1
			}
		case *Group:
			if verifyThreshold(t, signers) {
				signed += 1
			}
		default:
			panic("invalid group member type")
		}
	}
	return signed >= g.Threshold
}

func verifyGroupSignature(srvc *native.NativeService, g *Group, signers []Signer) bool {
	if !verifyThreshold(g, signers) {
		return false
	}

	for _, signer := range signers {
		key, err := encodeID(signer.id)
		if err != nil {
			return false
		}
		pk, err := getPk(srvc, key, signer.index)
		if err != nil {
			return false
		}
		if pk.revoked {
			return false
		}
		if checkWitness(srvc, pk.key) != nil {
			return false
		}
	}
	return true
}
