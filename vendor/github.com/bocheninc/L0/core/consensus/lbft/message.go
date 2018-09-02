// Copyright (C) 2017, Beijing Bochen Technology Co.,Ltd.  All rights reserved.
//
// This file is part of L0
//
// The L0 is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The L0 is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package lbft

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"

	"github.com/bocheninc/L0/components/utils"
	"github.com/bocheninc/L0/core/types"
)

//Request Define struct
type Request struct {
	ID     int64
	Time   uint32
	Height uint32
	Txs    types.Transactions
}

//Name key name
func (msg *Request) Name() string {
	keys := make([]string, 0)
	r := &Request{
		ID:   msg.ID,
		Time: msg.Time,
		//Height: msg.Height,
		Txs: msg.Txs,
	}
	hash := sha256.Sum256(utils.Serialize(r))
	keys = append(keys, hex.EncodeToString(hash[:]))
	keys = append(keys, strconv.Itoa(len(msg.Txs)))
	return strings.Join(keys, "-")
}

//PrePrepare Define struct
type PrePrepare struct {
	PrimaryID string
	SeqNo     uint32
	Height    uint32
	OptHash   string
	// Digest    string
	MerkleRoot string
	Quorum     int
	Request    *Request
	Chain      string
	ReplicaID  string
}

//Prepare Define struct
type Prepare struct {
	PrimaryID string
	SeqNo     uint32
	Height    uint32
	OptHash   string
	Digest    string
	Quorum    int
	Chain     string
	ReplicaID string
}

//Commit Define struct
type Commit struct {
	PrimaryID string
	SeqNo     uint32
	Height    uint32
	OptHash   string
	Digest    string
	Quorum    int
	Chain     string
	ReplicaID string
}

//Committed Define struct
type Committed struct {
	SeqNo     uint32
	Height    uint32
	Digest    string
	Txs       types.Transactions
	ErrTxs    types.Transactions
	Chain     string
	ReplicaID string
}

//FetchCommitted Define struct
type FetchCommitted struct {
	SeqNo     uint32
	Chain     string
	ReplicaID string
}

//ViewChange Define struct
type ViewChange struct {
	ID            string
	Priority      int64
	PrimaryID     string
	SeqNo         uint32
	Height        uint32
	OptHash       string
	LastPrimaryID string
	ReplicaID     string
	Chain         string
}

//MessageType
type MessageType uint32

const (
	MESSAGEUNDEFINED      MessageType = 0
	MESSAGEREQUEST        MessageType = 1
	MESSAGEPREPREPARE     MessageType = 2
	MESSAGEPREPARE        MessageType = 3
	MESSAGECOMMIT         MessageType = 4
	MESSAGECOMMITTED      MessageType = 5
	MESSAGEFETCHCOMMITTED MessageType = 6
	MESSAGEVIEWCHANGE     MessageType = 7
)

//Message Define lbft message struct
type Message struct {
	// Types that are valid to be assigned to Payload:
	//	*Request
	//	*PrePrepare
	//	*Prepare
	//	*Commit
	//	*Committed
	//	*FetchCommitted
	//	*ViewChange
	Type    MessageType
	Payload []byte
}

//GetRequestBatch
func (m *Message) GetRequest() *Request {
	if m.Type == MESSAGEREQUEST {
		x := &Request{}
		if err := utils.Deserialize(m.Payload, x); err != nil {
			panic(err)
		}
		return x
	}
	return nil
}

//GetPrePrepare
func (m *Message) GetPrePrepare() *PrePrepare {
	if m.Type == MESSAGEPREPREPARE {
		x := &PrePrepare{}
		if err := utils.Deserialize(m.Payload, x); err != nil {
			panic(err)
		}
		return x
	}
	return nil
}

//Get Prepare
func (m *Message) GetPrepare() *Prepare {
	if m.Type == MESSAGEPREPARE {
		x := &Prepare{}
		if err := utils.Deserialize(m.Payload, x); err != nil {
			panic(err)
		}
		return x
	}
	return nil
}

//GetCommit
func (m *Message) GetCommit() *Commit {
	if m.Type == MESSAGECOMMIT {
		x := &Commit{}
		if err := utils.Deserialize(m.Payload, x); err != nil {
			panic(err)
		}
		return x
	}
	return nil
}

//GetCommitted
func (m *Message) GetCommitted() *Committed {
	if m.Type == MESSAGECOMMITTED {
		x := &Committed{}
		if err := utils.Deserialize(m.Payload, x); err != nil {
			panic(err)
		}
		return x
	}
	return nil
}

//GetFetchCommitted
func (m *Message) GetFetchCommitted() *FetchCommitted {
	if m.Type == MESSAGEFETCHCOMMITTED {
		x := &FetchCommitted{}
		if err := utils.Deserialize(m.Payload, x); err != nil {
			panic(err)
		}
		return x
	}
	return nil
}

//GetViewChange
func (m *Message) GetViewChange() *ViewChange {
	if m.Type == MESSAGEVIEWCHANGE {
		x := &ViewChange{}
		if err := utils.Deserialize(m.Payload, x); err != nil {
			panic(err)
		}
		return x
	}
	return nil
}
