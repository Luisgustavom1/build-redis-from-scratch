package main

import (
	"bufio"
	"os"
	"sync"
	"time"

	"github.com/Luisgustavom1/build-redis-from-scratch/resp"
)

type Aof struct {
	file *os.File
	rd   *bufio.Reader
	mu   sync.Mutex
}

func NewAof(path string) (*Aof, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	aof := &Aof{
		file: f,
		mu:   sync.Mutex{},
		rd:   bufio.NewReader(f),
	}

	go func() {
		for {
			aof.mu.Lock()
			aof.file.Sync()
			aof.mu.Unlock()

			// every second sync the file
			time.Sleep(1 * time.Second)
		}
	}()

	return aof, nil
}

func (aof *Aof) Close() error {
	aof.mu.Lock()
	defer aof.mu.Unlock()
	return aof.file.Close()
}

func (aof *Aof) Write(value resp.Value) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	_, err := aof.file.Write(value.Marshal())
	if err != nil {
		return err
	}

	return nil
}

func (aof *Aof) Read(fn func(resp.Value)) {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	res := resp.NewResp(aof.rd)

	for {
		value, err := res.Read()
		if err != nil {
			return
		}
		fn(value)
	}
}
