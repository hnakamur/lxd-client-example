package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/lxc/lxd"
	"github.com/pkg/errors"
)

type nopWriteCloser struct {
	io.Writer
}

func NopWriteCloser(w io.Writer) nopWriteCloser {
	return nopWriteCloser{w}
}

func (nopWriteCloser) Close() error { return nil }

type stdWriter struct {
	w      *os.File
	header []byte
}

func newStdWriter(w *os.File, header []byte) stdWriter {
	return stdWriter{w: w, header: header}
}

func (w stdWriter) Write(p []byte) (int, error) {
	_, err := w.w.Write(w.header)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return w.w.Write(p)
}

func (w stdWriter) Close() error {
	fmt.Fprintf(os.Stderr, "Close called for %s\n", string(w.header))
	return nil
}

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
	stdout := newStdWriter(os.Stdout, []byte("out: "))
	stderr := newStdWriter(os.Stderr, []byte(color.RedString("err: ")))
	rc, err := client.Exec(container, cmd,
		env, stdin, stdout, stderr, nil, 0, 0)
	if err != nil {
		log.Fatalf("failed to create client")
	}
	fmt.Printf("rc=%d\n", rc)
}
