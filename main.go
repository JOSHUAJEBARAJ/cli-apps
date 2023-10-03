package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
)

func openFile(filename string) (*os.File, error) {

	newFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return newFile, nil

}
func WordCounter(r io.Reader) (int, error) {
	fileScanner := bufio.NewScanner(r)
	fileScanner.Split(bufio.ScanWords)
	count := 0
	for fileScanner.Scan() {
		count++
	}

	if err := fileScanner.Err(); err != nil {
		return count, err
	}
	return count, nil

}

func CharCounter(r io.ReadCloser) (int, error) {
	fileScanner := bufio.NewScanner(r)
	fileScanner.Split(bufio.ScanBytes)
	count := 0
	for fileScanner.Scan() {
		count++
	}

	if err := fileScanner.Err(); err != nil {
		return count, err
	}
	return count, nil

}
func LineCounter(r io.Reader) (int, error) {
	// create a buffer
	buf := make([]byte, 32*1024)
	count := 0
	// create a line seperator
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func main() {
	var lFlag = flag.Bool("l", false, "use to find the line count")
	var cFlag = flag.Bool("c", false, "use to find the character count")
	var wFlag = flag.Bool("w", false, "use to find the word count")
	flag.Usage = func() {

		fmt.Println("This program was created by Joshua Jebaraj")

		flag.PrintDefaults()
	}
	flag.Parse()
	if len(flag.Args()) <= 0 {
		flag.Usage()
		os.Exit(1)
	}
	filenames := flag.Args()
	//fmt.Println(flag.Args(), *lFlag, *cFlag, *wFlag)
	for _, filename := range filenames {
		file, err := openFile(filename)
		if err != nil {
			handleError(err)
		}
		if *lFlag {
			lc, err := LineCounter(file)
			if err != nil {
				handleError(err)
			}
			fmt.Println(lc, filename)
		}
		if *wFlag {
			file.Seek(0, 0)
			wc, err := WordCounter(file)
			if err != nil {
				handleError(err)
			}
			fmt.Println(wc, filename)
		}
		if *cFlag {
			file.Seek(0, 0)
			cc, err := CharCounter(file)
			if err != nil {
				handleError(err)
			}
			fmt.Println(cc, filename)
		}
		if !*cFlag && !*wFlag && !*lFlag {
			file.Seek(0, 0)
			lc, err := LineCounter(file)
			if err != nil {
				handleError(err)
			}
			file.Seek(0, 0)
			wc, err := WordCounter(file)
			if err != nil {
				handleError(err)
			}
			file.Seek(0, 0)
			cc, err := CharCounter(file)
			if err != nil {
				handleError(err)
			}
			fmt.Println(lc, wc, cc, filename)
		}

	}
}
func handleError(er error) {
	fmt.Fprintln(os.Stderr, er)
	os.Exit(1)
}
