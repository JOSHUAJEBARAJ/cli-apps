package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type config struct {
	search_string  string
	filename       string
	output         string
	case_sensitive bool
}

type result struct {
	lines             []string
	filename          string
	printwithfilename bool
}

func parseFlags() config {
	o := flag.String("o", "", "To send the output to the file")
	i := flag.Bool("i", false, "case sensitive")
	flag.Parse()
	if len(flag.Args()) < 1 {
		flag.PrintDefaults()
		os.Exit(1)
	}
	input_search_text := flag.Arg(0)
	filename := flag.Arg(1)
	c := config{search_string: input_search_text, filename: filename, output: *o, case_sensitive: *i}
	return c
}

func openFile(c config) (*os.File, error) {
	if c.filename == "" {
		return os.Stdin, nil
	}
	f, err := os.Open(c.filename)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func search(r io.Reader, c config, wg *sync.WaitGroup, resultsChan chan<- result, mu *sync.Mutex) {
	defer wg.Done()
	scanner := bufio.NewScanner(r)
	lines := []string{}
	for scanner.Scan() {
		if c.case_sensitive {
			if strings.Contains(strings.ToLower(scanner.Text()), strings.ToLower(c.search_string)) {
				lines = append(lines, scanner.Text())
			}
		} else {
			if strings.Contains(scanner.Text(), c.search_string) {
				lines = append(lines, scanner.Text())
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error searching %s: %v\n", c.filename, err)
		return
	}

	// Protect the critical section with a mutex lock
	mu.Lock()
	resultsChan <- result{lines: lines, filename: c.filename, printwithfilename: true}
	mu.Unlock()
}

func walk(path string) (allFiles []string) {
	filepath.Walk(path, func(file string, info os.FileInfo, err error) error {
		// skip the root
		if info.IsDir() {
			return nil
		}
		if file != path {
			allFiles = append(allFiles, file)
		}
		return nil
	})
	return allFiles
}
func printresult(c config, r result) {
	if c.output == "" {
		for _, line := range r.lines {
			if r.printwithfilename {
				fmt.Printf("%s : %s\n", r.filename, line)
			} else {
				fmt.Println(line)
			}
		}
		return
	}
	outputfilename := c.output
	outputFile, err := os.Create(outputfilename)
	if err != nil {
		handleError(err)
	}
	defer outputFile.Close()
	writer := bufio.NewWriter(outputFile)
	for _, line := range r.lines {
		if r.printwithfilename {
			_, err := writer.WriteString(r.filename + ":" + line + "\n")
			if err != nil {
				handleError(err)
			}
		} else {
			_, err := writer.WriteString(line + "\n")
			if err != nil {
				handleError(err)
			}
		}
	}
	err = writer.Flush()
	if err != nil {
		handleError(err)
	}
}

func main() {
	c := parseFlags()
	var wg sync.WaitGroup
	resultsChan := make(chan result, 100) // Buffered channel to hold results
	var mu sync.Mutex

	// Check if the filename is a directory
	info, err := os.Stat(c.filename)
	if err != nil {
		handleError(err)
	}

	if info.IsDir() {
		files := walk(c.filename)
		for _, file := range files {
			c.filename = file
			f, err := openFile(c)
			if err != nil {
				handleError(err)
			}
			wg.Add(1)
			go search(f, c, &wg, resultsChan, &mu)
		}
	} else {
		wg.Add(1)
		f, err := openFile(c)
		if err != nil {
			handleError(err)
		}
		go search(f, c, &wg, resultsChan, &mu)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	for res := range resultsChan {
		printresult(c, res)
	}
}

func handleError(err error) {
	fmt.Fprintln(os.Stderr, "Error:", err)
	os.Exit(1)
}
