package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"

	"github.com/lxc/lxd"
)

type nopWriteCloser struct {
	io.Writer
}

func NopWriteCloser(w io.Writer) nopWriteCloser {
	return nopWriteCloser{w}
}

func (nopWriteCloser) Close() error { return nil }

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	client, err := lxd.NewClient(&lxd.DefaultConfig, "local")
	if err != nil {
		log.Fatalf("failed to create client")
	}

	var env map[string]string
	stdin := ioutil.NopCloser(new(bytes.Buffer))
	outBuf := new(bytes.Buffer)
	stdout := NopWriteCloser(outBuf)
	errBuf := new(bytes.Buffer)
	stderr := NopWriteCloser(errBuf)
	rc, err := client.Exec("centos7", []string{"hostname"}, env,
		stdin, stdout, stderr, nil, 0, 0)
	if err != nil {
		log.Fatalf("failed to create client")
	}

	log.Printf("rc=%d, stdout=%q, stderr=%q", rc, outBuf.String(), errBuf.String())
}
