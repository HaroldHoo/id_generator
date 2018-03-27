/**
 * Copyright 2018 harold. All rights reserved.
 * Author: harold
 * Mail: mail@yaolong.me
 * Date: 2018-03-22
 */

package id_generator

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

const (
	MAX_NEXT_ID     = 1 << 14
	MAX_DATA_ID     = 1 << 10
	MAX_INSTANCE_ID = 1 << 8
)

var (
	currentTimestamp  uint64
	currentNextId     uint64
	DefaultInstanceId uint64
	DefaultCacheFile  string

	ErrNextIdOutOf     = errors.New(fmt.Sprintf("nextID out of %d", MAX_NEXT_ID))
	ErrDataIdOutOf     = errors.New(fmt.Sprintf("dataId out of %d", MAX_DATA_ID))
	ErrInstanceIdOutOf = errors.New(fmt.Sprintf("instanceId out of %d", MAX_INSTANCE_ID))
)

type id_generator struct {
	timestamp uint64
	extraId   struct {
		instanceId uint64
		dataId     uint64
		nextId     uint64
	}
	mu     sync.Mutex
	result uint64
	flock *syscall.Flock_t
}

func New() (id *id_generator) {
	id = &id_generator{}
	id.flock = &syscall.Flock_t{
		Type:  syscall.F_WRLCK,
	}
	return
}

var g_id_generator *id_generator

func Init() *id_generator {
	if g_id_generator != nil {
		return g_id_generator
	}
	g_id_generator = New()
	return g_id_generator
}

func NextId(dataId uint64) (ret uint64, err error) {
	id_gen := Init()
	return id_gen.NextId(dataId)
}

func SetDefaultCacheFile(f string) {
	if DefaultCacheFile != "" {
		// panic("The id_generator.DefaultCacheFile has been set.")
		return
	}

	DefaultCacheFile = f
}

var cacheFileHandler *os.File

func cacheFile() *os.File {
	if cacheFileHandler == nil {
		var err error
		cacheFileHandler, err = os.Create(DefaultCacheFile)
		if err != nil {
			panic(err)
		}
	}
	return cacheFileHandler
}

func (id *id_generator) setTimestampCache(t uint64, u uint64) {
	if DefaultCacheFile == "" {
		currentTimestamp = t
		currentNextId = u
		return
	}

	f := cacheFile()
	syscall.FcntlFlock(uintptr(f.Fd()), syscall.F_SETLKW, id.flock)
	defer syscall.FcntlFlock(uintptr(f.Fd()), syscall.F_UNLCK, id.flock)

	f.Seek(0, os.SEEK_SET)
	f.WriteString(fmt.Sprintf("%d", (t<<14 | u)))
	f.Sync()

	return
}

func (id *id_generator) getTimestampCache() (t uint64, n uint64) {
	if DefaultCacheFile == "" {
		return currentTimestamp, currentNextId
	}

	f := cacheFile()
	syscall.FcntlFlock(uintptr(f.Fd()), syscall.F_SETLKW, id.flock)
	defer syscall.FcntlFlock(uintptr(f.Fd()), syscall.F_UNLCK, id.flock)

	f.Seek(0, os.SEEK_SET)
	b := make([]byte, 46)
	ns, err := f.Read(b)
	if ns == 0 {
		return 0, 0
	}
	if err != nil {
		panic(err)
	}

	var num int64
	num, err = strconv.ParseInt(strings.TrimRight(string(b), "\x00"), 10, 64)
	if err != nil {
		panic(err)
	}

	t = uint64(num) >> 14
	n = uint64(num) & 16383

	return
}

func (id *id_generator) NextId(dataId uint64) (ret uint64, err error) {
	id.mu.Lock()
	defer id.mu.Unlock()

	id.result = 0
	id.timestamp = uint64(time.Now().Unix())

	t, n := id.getTimestampCache()
	if t == id.timestamp {
		n++
		id.setTimestampCache(id.timestamp, n)
		id.extraId.nextId = n
	} else {
		id.setTimestampCache(id.timestamp, 0)
		id.extraId.nextId = 0
	}

	id.extraId.instanceId = DefaultInstanceId
	id.extraId.dataId = dataId

	switch true {
	case id.extraId.nextId >= MAX_NEXT_ID:
		err = ErrNextIdOutOf
		return
	case id.extraId.dataId >= MAX_DATA_ID:
		err = ErrDataIdOutOf
		return
	case id.extraId.instanceId >= MAX_INSTANCE_ID:
		err = ErrInstanceIdOutOf
		return
	}

	id.result |= id.timestamp << 32          // 32bit
	id.result |= id.extraId.instanceId << 24 // 08bit
	id.result |= id.extraId.dataId << 14     // 10bit
	id.result |= id.extraId.nextId           // 14bit
	return id.result, nil
}

func GetTimestamp(id uint64) uint64 {
	return id >> 32
}

func GetInstanceId(id uint64) uint64 {
	return id >> 24 & 255
}

func GetDataId(id uint64) uint64 {
	return id >> 14 & 1023
}

func GetNextId(id uint64) uint64 {
	return id & 16383
}

func (id *id_generator) DebugPrint() {
	fmt.Printf("%s\n", "DEBUG--------------------")

	var s string
	s = fmt.Sprintf("%b", id.result)
	fmt.Printf("%s%s\n", strings.Repeat("0", 64-len([]rune(s))), s)

	s = fmt.Sprintf("%b", id.timestamp)
	fmt.Printf("%s%s\n", strings.Repeat("0", 32-len([]rune(s))), s)

	s = fmt.Sprintf("%b", id.extraId.instanceId)
	fmt.Printf("instanceId:%s%s%s\n", strings.Repeat(" ", 32-11), strings.Repeat("0", 8-len([]rune(s))), s)

	s = fmt.Sprintf("%b", id.extraId.dataId)
	fmt.Printf("dataId:%s%s%s\n", strings.Repeat(" ", 32+8-7), strings.Repeat("0", 10-len([]rune(s))), s)

	s = fmt.Sprintf("%b", id.extraId.nextId)
	fmt.Printf("nextId:%s%s%s\n", strings.Repeat(" ", 32+8+10-7), strings.Repeat("0", 14-len([]rune(s))), s)

	fmt.Printf("%s\n\n", "DEBUG--------------------")
}
