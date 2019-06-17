package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hallazzang/syso"
	"github.com/hallazzang/syso/pkg/coff"
)

var (
	configFile string
	outFile    string
)

func init() {
	flag.StringVar(&configFile, "c", "syso.json", "config file name")
	flag.StringVar(&outFile, "o", "out.syso", "output file name")
	flag.Parse()
}

func main() {
	fcfg, err := os.Open(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open config file: %v\n", err)
		os.Exit(1)
	}
	defer fcfg.Close()

	cfg, err := syso.ParseConfig(fcfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse config: %v\n", err)
		os.Exit(1)
	}

	c := coff.New()
	for i, icon := range cfg.Icons {
		if err := syso.EmbedIcon(c, icon); err != nil {
			fmt.Fprintf(os.Stderr, "failed to embed icon #%d: %v\n", i, err)
			os.Exit(1)
		}
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
