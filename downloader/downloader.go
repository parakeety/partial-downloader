package downloader

import (
	"errors"
	"io/ioutil"
	"net/http"
	"runtime"
	"strconv"

	"golang.org/x/sync/errgroup"
)

type Downloader struct {
	URL      string
	Filename string
}

func (d *Downloader) Start() error {
	size, err := d.head()
	if err != nil {
		return err
	}

	numCPU := runtime.NumCPU()
	unit := int(size) / numCPU

	data := newData(numCPU)
	var eg errgroup.Group
	for i := 0; i < numCPU; i++ {
		index := i
		start := unit * i
		end := unit*(i+1) - 1
		if size < end {
			end = size
		}
		eg.Go(func() error {
			return d.rangeAccess(data, index, start, end)
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	if err := data.write(d.Filename); err != nil {
		return err
	}

	return nil
}

func (d *Downloader) head() (int, error) {
	headResp, err := http.Head(d.URL)
	if err != nil {
		return 0, err
	}

	if headResp.Header.Get("Accept-Ranges") != "bytes" {
		return 0, errors.New("specified url doesn't support accept-ranges")
	}

	return int(headResp.ContentLength), nil
}

func (d *Downloader) rangeAccess(data *Data, index, start, end int) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", d.URL, nil)
	if err != nil {
		return err
	}

	rangeBytes := "bytes=" + strconv.Itoa(start) + "-" + strconv.Itoa(end)
	req.Header.Add("Range", rangeBytes)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	p, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	data.setPartialData(index, p)

	return nil
}
