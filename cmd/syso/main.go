package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hallazzang/syso/pkg/coff"
	"github.com/hallazzang/syso/pkg/ico"
	"github.com/hallazzang/syso/pkg/rsrc"
)

var (
	icoFile string
	outFile string
)

func init() {
	flag.StringVar(&icoFile, "ico", "", ".ico file to embed")
	flag.StringVar(&outFile, "o", "out.syso", "output file name")
	flag.Parse()
}

func main() {
	if icoFile == "" {
		fmt.Fprintln(os.Stderr, "icon file is not provided")
		os.Exit(1)
	}

	fico, err := os.Open(icoFile)
	if err != nil {
		panic(err)
	}
	defer fico.Close()

	ig, err := ico.DecodeAll(fico)
	if err != nil {
		panic(err)
	}
	for i, img := range ig.Images {
		img.ID = i + 100
	}

	c := coff.New()
	r := rsrc.New()
	if err := r.AddIconsByID(1, ig); err != nil {
		panic(err)
	}
	if err := c.AddSection(r); err != nil {
		panic(err)
	}

	fout, err := os.Create(outFile)
	if err != nil {
		panic(err)
	}
	defer fout.Close()

	if _, err := c.WriteTo(fout); err != nil {
		panic(err)
	}

	fmt.Printf("successfully generated syso file to %s", outFile)
}
