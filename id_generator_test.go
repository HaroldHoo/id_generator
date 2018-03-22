/**
 * Copyright 2018 harold. All rights reserved.
 * Filename: ./id_generator/id_generator_test.go
 * Author: harold
 * Mail: mail@yaolong.me
 * Date: 2018-03-22
 */

package id_generator

import (
	"fmt"
	"testing"
	"time"
)

func TestDebugPrint(t *testing.T) {
	fmt.Printf("\n----------------- TestDebugPrint\n")
	DefaultInstanceId = 123
	id := New()
	id.NextId(1022)
	id.DebugPrint()
}

func TestErrInstanceIdOutOf(t *testing.T) {
	fmt.Printf("\n----------------- TestErrInstanceIdOutOf\n")
	DefaultInstanceId = 256
	id := New()
	_, e := id.NextId(1022)
	if e != ErrInstanceIdOutOf {
		t.Fail()
	}
	fmt.Println("OK")
}

func TestErrDataIdOutOf(t *testing.T) {
	fmt.Printf("\n----------------- TestErrDataIdOutOf\n")
	DefaultInstanceId = 255
	id := New()
	_, e := id.NextId(1024)
	if e != ErrDataIdOutOf {
		t.Fail()
	}
	fmt.Println("OK")
}

func TestErrNextIdOutOf(t *testing.T) {
	fmt.Printf("\n----------------- TestErrNextIdOutOf\n")
	DefaultInstanceId = 255
	id := New()

	var e error
	for i := 0; i < 16384; i++ {
		_, e = id.NextId(1022)
	}
	if e != ErrNextIdOutOf {
		t.Fail()
	}
	fmt.Println("OK")
}

func TestGenerator(t *testing.T) {
	fmt.Printf("\n----------------- TestGenerator\n")
	id := Init()

	var (
		gid     uint64
		e       error
		removed = make(map[uint64]int)
	)
	for i := 0; i < 16383; i++ {
		gid, e = id.NextId(1022)
		if e != nil {
			t.Fatal(e)
		}
		removed[gid] = 0
	}
	if len(removed) != 16383 {
		t.Fail()
	}

	var t1, t2, t3, t4, n1, n2, n3, n4 uint64
	DefaultInstanceId = 212
	for {
		if time.Now().Nanosecond()/1e6 == 100 {
			id.NextId(101)
			t1 = id.timestamp
			n1 = id.extraId.nextId
			id.NextId(11)
			t2 = id.timestamp
			n2 = id.extraId.nextId
			break
		}
	}
	time.Sleep(1 * time.Second)
	for {
		if time.Now().Nanosecond()/1e6 == 100 {
			id.NextId(12)
			t3 = id.timestamp
			n3 = id.extraId.nextId
			id.NextId(13)
			t4 = id.timestamp
			n4 = id.extraId.nextId
			break
		}
	}

	if t1 != t2 || t3 != t4 {
		t.Errorf("%s\n", "timestamp not equle")
	}
	if n1 != n3 || n2 != n4 {
		t.Errorf("%s\n", "nextId not init")
	}
	if n1 == n2 || n3 == n4 {
		t.Errorf("%s\n", "nextId not unique at same time")
	}

	fmt.Println("OK")
}

