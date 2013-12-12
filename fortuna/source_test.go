package fortuna

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func newSourceChannel(cs *SourceChannel) error {
	f, err := os.Open("/dev/random")
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	io.CopyN(buf, f, 4096)
	f.Close()
	go func() {
		for {
			if buf.Len() == 0 {
				return
			}
			e := make([]byte, 24)
			_, err := buf.Read(e)
			if err != nil {
				fmt.Fprintf(os.Stderr, "fortuna: error on buffer (%v)\n", err)
				return

			}
			cs.In <- e
			<-time.After(30 * time.Millisecond)
		}
	}()
	return nil
}

func TestSourceChannel(t *testing.T) {
	rng := New(nil)
	cs := NewSourceChannel(rng, 1)
	cs.Start(4)
	err := newSourceChannel(cs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		t.FailNow()
	}
	err = newSourceChannel(cs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		t.FailNow()
	}
	err = newSourceChannel(cs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		t.FailNow()
	}
	err = newSourceChannel(cs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		t.FailNow()
	}

	<-time.After(1 * time.Second)
	var p = make([]byte, 16384)
	n, err := rng.Read(p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		t.FailNow()
	} else if n != 16384 {
		fmt.Fprintf(os.Stderr, "fortuna: read %d bytes, expected 16384\n", n)
		t.FailNow()
	}
	ioutil.WriteFile("/tmp/fortuna.out", p, 0644)
}

func TestSourceWriter(t *testing.T) {
	rng := New(nil)
	sw := NewSourceWriter(rng, 2)
	f, err := os.Open("/dev/random")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		t.FailNow()
	}
	buf := new(bytes.Buffer)
	io.CopyN(buf, f, 4096)
	f.Close()

	n, err := io.Copy(sw, buf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		t.FailNow()
	} else if n != 4096 {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		t.FailNow()
	}
}