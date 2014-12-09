package main

import (
	"crypto/sha1"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var verbose = flag.Bool("v", false, "true to print extra info")

func main() {
	log.SetFlags(0)
	flag.Parse()

	files := flag.Args()
	sums := map[[sha1.Size]byte]string{}
	hasher := sha1.New()
	var sum [sha1.Size]byte
	ntot := int64(0)

	if len(files) == 0 {
		data, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		files = strings.Split(string(data), "\n")
	}

	for i, fname := range files {
		fi, err := os.Stat(fname)
		if err != nil {
			log.Print(err)
			continue
		}
		if fi.IsDir() {
			continue
		}

		f, err := os.Open(fname)
		if err != nil {
			log.Println(err)
			continue
		}

		hasher.Reset()
		n, err := io.Copy(hasher, f)
		if err != nil {
			log.Println(err)
			f.Close()
			continue
		}
		ntot += n

		if *verbose {
			if i%(len(files)/100) == 0 {
				progress := i * 100 / len(files)
				fmt.Printf("[INFO] %v%% done - checked %v/%v files (%v MB)\n", progress, i, len(files), ntot/1024/1024)
			}
		}

		data := hasher.Sum(nil)
		copy(sum[:], data)

		if prev, ok := sums[sum]; ok {
			if *verbose {
				fmt.Printf(" [DUP] %v duplicates %v\n", fname, prev)
			} else {
				fmt.Println(fname)
			}
		}
		sums[sum] = fname

		f.Close()
	}
}
