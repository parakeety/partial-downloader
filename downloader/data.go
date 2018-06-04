package downloader

import (
	"io/ioutil"
	"sync"
)

type data struct {
	data [][]byte
	mu   sync.Mutex
}

type partialData struct {
	data  []byte
	index int
}

func newData(count int) *data {
	return &data{
		data: make([][]byte, count),
	}
}

func (d *data) setPartialData(p *partialData) {
	d.mu.Lock()
	d.data[p.index] = p.data
	d.mu.Unlock()
}

func (d *data) write(filename string) error {
	var sum []byte
	for _, p := range d.data {
		sum = append(sum, p...)
	}

	return ioutil.WriteFile(filename, sum, 0666)
}
