package modulo37

import (
	"bytes"
	"fmt"
	"github.com/iwdgo/testingfiles"
	"io"
	"os"
	"testing"
	"time"
)

/*

About testing

Example would be like Test...
Log format, if used, and includes date, time...
pfile := iotest.NewWriteLogger(t.Name(),os.Stdout) and no valid reference can be created

Go test // Output: does not make a difference between one crlf and several.

To handle output, you can write the func with a io.Writer parameter like below but it requires to update or to write
code with test in mind. This is not required by Go.

func modulo37(f io.Writer) {
	for i := 1; i < 100; i++ {
		if i%3 == 0 {
			fmt.Fprint(f,"Open")
		}
		if i%7 == 0 {
			fmt.Fprint(f,"Source")
		}
		if (i%3 != 0) && (i%7 != 0) {
			fmt.Fprintf(f,"%d\n",i)
		} else {
			fmt.Fprintln(f)
		}
	}
}

You can conceive your test to pass the produced file and all output of the func is like fmt.Fprintf(pfile,...)

func TestModulo37(t *testing.T) {
	prodFileName := "moduloprod.txt"
	pfile, err := os.Create(prodFileName)
	defer pfile.Close()
	if err != nil {
		t.Errorf("Produced file creation %s failed with %v",prodFileName,err)
	}


	modulo37(pfile) // t.Log is using stdErr and looks confusing
	pfile.Close()

	testingfiles.FileCompare(t,"moduloref.txt",prodFileName)
}

*/

// Test is piping output to a file which is checked against reference.
func TestModulo37(t *testing.T) {
	testingfiles.OutputDir("output")
	prodFileName := "moduloprod.txt"
	pfile, err := os.Create(prodFileName)
	if err != nil {
		t.Errorf("produced file creation %s failed with %v", prodFileName, err)
	}

	// Capture stdout.
	stdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "modulo : piping failed with ", err)
		os.Exit(1)
	}
	os.Stdout = w
	outC := make(chan []byte)
	go func() {
		var buf bytes.Buffer
		_, err := io.Copy(&buf, r)
		r.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "testing: copying pipe: %v\n", err)
			os.Exit(1)
		}
		outC <- buf.Bytes() //.String()
	}()

	start := time.Now()
	ok := true

	/* Clean up in a deferred call so we can recover if the example panics. */
	defer func() {
		err := recover()
		if err != nil { // If here because of panic
			t.Error(err)
			panic(err) // Testing output has no value
		}

		dstr := fmt.Sprintf("%.4fs", time.Since(start).Seconds())

		// Close pipe, restore stdout, get output.
		w.Close()
		os.Stdout = stdout // Restoring Stdout
		out := <-outC
		pfile.Write(out)
		pfile.Close()

		if err = testingfiles.FileCompare(prodFileName, "modulowant.txt"); err != nil {
			t.Errorf("%s : %v\n", dstr, err)
		}

		ok = err == nil
	}()

	if !ok {
		t.Errorf("Opening pipe failed with %v", err)
	}
	modulo37()
	// All output handling is in defer

}

/* */
func BenchmarkModulo37(b *testing.B) {
	prodFileName := "modulobench.txt"
	pfile, err := os.Create(prodFileName)
	if err != nil {
		b.Errorf("produced file creation %s failed with %v", prodFileName, err)
	}

	// Capture stdout.
	stdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "modulo : piping failed with ", err)
		os.Exit(1)
	}
	os.Stdout = w
	outC := make(chan []byte)
	go func() {
		var buf bytes.Buffer
		_, err := io.Copy(&buf, r)
		r.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "testing: copying pipe: %v\n", err)
			os.Exit(1)
		}
		outC <- buf.Bytes() //.String()
	}()

	/* Clean up in a deferred call so we can recover if the example panics. */
	defer func() {
		err := recover()
		if err != nil { // If here because of panic
			b.Error(err)
			panic(err) // Testing output has no value
		}

		// Close pipe, restore stdout, get output.
		w.Close()
		os.Stdout = stdout // Restoring Stdout
		out := <-outC
		pfile.Write(out)
		pfile.Close()
	}()

	// run the function b.N times
	for n := 0; n < b.N; n++ {
		modulo37()
	}
}
