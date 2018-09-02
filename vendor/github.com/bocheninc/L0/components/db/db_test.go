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

package db

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"testing"
)

func TestReadAndWrite(t *testing.T) {
	db := NewDB(DefaultConfig())
	defer os.RemoveAll(config.DbPath)
	// Put
	err := db.Put("default", []byte("foo"), []byte("bar"))
	if err != nil {
		t.Fatalf("faild to put, err: [%s]", err)
	}
	// Get
	value, err1 := db.Get("default", []byte("foo"))
	if err1 != nil {
		t.Fatalf("faild to get, err: [%s]", err1)
	}
	if !bytes.Equal(value, []byte("bar")) {
		t.Fatal("value not equal")
	}
}

func TestDelete(t *testing.T) {

	db := NewDB(DefaultConfig())
	defer os.RemoveAll(config.DbPath)

	err := db.Put("default", []byte("foo"), []byte("bar"))
	if err != nil {
		t.Fatalf("faild to put, err: [%s]", err)
	}
	db.Delete("default", []byte("foo"))
	value, err1 := db.Get("default", []byte("foo"))
	if err1 != nil {
		t.Fatalf("faild to delete, err: [%s]", err)
	}
	if value != nil {
		t.Fatalf("faild to put")
	}
}

func TestGetByRangeOrPrefix(t *testing.T) {

	db := NewDB(DefaultConfig())
	defer os.RemoveAll(config.DbPath)
	for i := 0; i < 10; i++ {
		key := []byte("key_" + strconv.Itoa(i))
		value := []byte("value_" + strconv.Itoa(i))

		err := db.Put("balance", key, value)
		if err != nil {
			t.Fatalf("faild to put, err: [%s]", err)
		}
		key1 := []byte("key_1" + strconv.Itoa(i))
		value1 := []byte("value_1" + strconv.Itoa(i))

		err = db.Put("balance", key1, value1)
		if err != nil {
			t.Fatalf("faild to put, err: [%s]", err)
		}
	}

	values := db.GetByRange("balance", []byte("key_1"), []byte("key_3"))

	for _, v := range values {
		fmt.Println("key: ", string(v.Key), "value: ", string(v.Value))
	}

	fmt.Println("-------------------------------")

	values1 := db.GetByPrefix("balance", []byte("key_1"))

	for _, v := range values1 {
		fmt.Println("key: ", string(v.Key), "value: ", string(v.Value))
	}
}

func TestWriteBatch(t *testing.T) {
	db := NewDB(DefaultConfig())
	defer os.RemoveAll(config.DbPath)
	var writeBatchs []*WriteBatch

	for i := 0; i < 100; i++ {
		writeBatchs = append(writeBatchs, NewWriteBatch("balance", OperationPut, []byte("key"+strconv.Itoa(i)), []byte("value"+strconv.Itoa(i)), "balance"))
	}
	fmt.Println("start writeBatch...")

	var cnt int
	for i := 0; i < 100; i++ {
		fmt.Println("times: ", cnt)
		db.AtomicWrite(writeBatchs)
		cnt++
	}

}
