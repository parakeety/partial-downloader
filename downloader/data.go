package downloader

import (
	"io/ioutil"
	"sync"
)

type Data struct {
	data [][]byte
	mu   sync.Mutex
}

func newData(count int) *Data {
	return &Data{
		data: make([][]byte, count),
	}
}

func (d *Data) setPartialData(index int, p []byte) {
	d.mu.Lock()
	d.data[index] = p
	d.mu.Unlock()
}

func (d *Data) write(filename string) error {
	var sum []byte
	for _, p := range d.data {
		sum = append(sum, p...)
	}

	return ioutil.WriteFile(filename, sum, 0666)
}
