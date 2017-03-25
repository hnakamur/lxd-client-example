package main

import (
	"bytes"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"

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
	flag.Parse()
	if flag.NArg() < 2 {
		flag.Usage()
		os.Exit(1)
	}
	args := flag.Args()

	config := &lxd.DefaultConfig
	remote, container := config.ParseRemoteAndContainer(args[0])
	cmd := args[1:]

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	client, err := lxd.NewClient(config, remote)
	if err != nil {
		log.Fatalf("failed to create client")
	}

	var env map[string]string
	stdin := ioutil.NopCloser(new(bytes.Buffer))
	outBuf := new(bytes.Buffer)
	stdout := NopWriteCloser(outBuf)
	errBuf := new(bytes.Buffer)
	stderr := NopWriteCloser(errBuf)
	rc, err := client.Exec(container, cmd,
		env, stdin, stdout, stderr, nil, 0, 0)
	if err != nil {
		log.Fatalf("failed to create client")
	}

	log.Printf("rc=%d, stdout=%q, stderr=%q", rc, outBuf.String(), errBuf.String())
}
