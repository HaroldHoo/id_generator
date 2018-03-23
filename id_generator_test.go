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
	"sync"
)

func TestDebugPrint(t *testing.T) {
	// DefaultCacheFile = "/tmp/id_g_test"
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
	_, e := id.NextId(1023)
	if e != ErrInstanceIdOutOf {
		t.Fatal()
	}
	fmt.Println("OK")
}

func TestErrDataIdOutOf(t *testing.T) {
	fmt.Printf("\n----------------- TestErrDataIdOutOf\n")
	DefaultInstanceId = 255
	id := New()
	_, e := id.NextId(1024)
	if e != ErrDataIdOutOf {
		t.Fatal()
	}
	fmt.Println("OK")
}

func TestErrNextIdOutOf(t *testing.T) {
	time.Sleep(1 * time.Second)
	fmt.Printf("\n----------------- TestErrNextIdOutOf\n")
	for {
		if time.Now().Nanosecond()/1e6 == 20 {
			DefaultInstanceId = 255
			id := New()

			if DefaultCacheFile != "" {
				setTimestampCache(uint64(time.Now().Unix()), 8000)
			}
			var e error
			var tmp uint64
			for i := 0; i < 16385; i++ {
				tmp, e = id.NextId(1022)
			}
			if e != ErrNextIdOutOf {
				t.Fatalf("\nmax next id: %v\n", GetNextId(tmp))
			}
			fmt.Println("OK")
			break
		}
	}
}

func TestGoroutine(t *testing.T){
	time.Sleep(1 * time.Second)
	fmt.Printf("\n----------------- TestGoroutine\n")
	m := make(map[uint64]int)
	var (
		mu sync.Mutex
	)
	c1 := make(chan int)
	c2 := make(chan int)
	go func(c chan int){
		for i := 0; i < 8192; i++ {
			id,e := NextId(255)
			if e != nil {
				fmt.Println(e)
				break
			}
			mu.Lock()
			m[id] = 0
			mu.Unlock()
		}
		c<-1
	}(c1)
	go func(c chan int){
		for i := 0; i < 8192; i++ {
			id,e := NextId(255)
			if e != nil {
				fmt.Println(e)
				break
			}
			mu.Lock()
			m[id] = 0
			mu.Unlock()
		}
		c<-1
	}(c2)

	<-c1
	<-c2

	if len(m) != 16384 {
		t.Fatalf("%s(%d)\n", "goroutine not safe", len(m))
	}

	fmt.Println("OK")
}

func TestGenerator(t *testing.T) {
	time.Sleep(1 * time.Second)
	fmt.Printf("\n----------------- TestGenerator\n")
	id := Init()

	var (
		gid     uint64
		e       error
		removed = make(map[uint64]int)
	)
	for i := 0; i < 16384; i++ {
		gid, e = id.NextId(1022)
		if e != nil {
			t.Fatal(e)
		}
		removed[gid] = 0
	}
	if len(removed) != 16384 {
		t.Fatal()
	}

	var t1, t2, t3, t4, n1, n2, n3, n4 uint64
	DefaultInstanceId = 212
	for {
		if time.Now().Nanosecond()/1e6 == 20 {
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
		if time.Now().Nanosecond()/1e6 == 20 {
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
		t.Fatalf("%s\n", "timestamp not equle")
	}
	if n1 != n3 || n2 != n4 {
		t.Fatalf("%s(%d,%d,%d,%d)\n", "nextId not init",n1,n3,n2,n4)
	}
	if n1 == n2 || n3 == n4 {
		t.Fatalf("%s\n", "nextId not unique at same time")
	}

	fmt.Println("OK")
}

func BenchmarkViaFile(b *testing.B){
	DefaultCacheFile = "/tmp/id_g_test"
	for i:=0; i<b.N; i++ {
		NextId(1)
	}
}

func BenchmarkViaMem(b *testing.B){
	DefaultCacheFile = ""
	for i:=0; i<b.N; i++ {
		NextId(1)
	}
}

