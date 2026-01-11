package utils

import "io"

type multiReadCloser struct {
	readers []io.ReadCloser
	current int
}

func MultiReadCloser(readers ...io.ReadCloser) io.ReadCloser {
	return &multiReadCloser{readers: readers}
}

func (c *multiReadCloser) Read(p []byte) (int, error) {
	for c.current < len(c.readers) {
		n, err := c.readers[c.current].Read(p)
		if err == io.EOF {
			c.readers[c.current].Close()
			c.current++
			continue
		}
		return n, err
	}
	return 0, io.EOF
}

func (c *multiReadCloser) Close() error {
	var firstErr error
	for i := c.current; i < len(c.readers); i++ {
		if err := c.readers[i].Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
