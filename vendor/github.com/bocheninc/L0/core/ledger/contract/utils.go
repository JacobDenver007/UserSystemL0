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

package contract

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"

	"github.com/bocheninc/L0/components/utils"
)

const (
	stringType = iota
)

var stateKeyDelimiter = "|"

func EnSmartContractKey(scAddr string, key string) string {
	return strings.Join([]string{scAddr, key}, stateKeyDelimiter)

}

func DeSmartContractKey(deKey string) (string, string) {
	split := strings.SplitN(deKey, stateKeyDelimiter, 2)

	return split[0], split[1]
}

func DoContractStateData(src []byte) ([]byte, error) {
	if len(src) == 0 {
		return nil, nil
	}

	buf := bytes.NewBuffer(src)
	tp, err := buf.ReadByte()
	if err != nil {
		return nil, err
	}

	switch tp {
	case stringType:
		len, err := utils.ReadVarInt(buf)
		if err != nil {
			return nil, err
		}
		data := make([]byte, len)
		buf.Read(data)
		return data, nil
	default:
		return nil, errors.New("not support states")
	}
}

func ConcrateStateJson(v interface{}) (*bytes.Buffer, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return &bytes.Buffer{}, err
	}

	buf := new(bytes.Buffer)
	buf.WriteByte(stringType)
	lenByte := utils.VarInt(uint64(len(data)))
	buf.Write(lenByte)
	buf.Write(data)

	return buf, nil
}

// func AppendContractPrefix(v []byte) (*bytes.Buffer, error) {
// 	data, err := json.Marshal(v)
// 	if err != nil {
// 		return &bytes.Buffer{}, err
// 	}

// 	buf := new(bytes.Buffer)
// 	buf.WriteByte(stringType)
// 	lenByte := utils.VarInt(uint64(len(data)))
// 	buf.Write(lenByte)
// 	buf.Write(data)

// 	return buf, nil
// }
