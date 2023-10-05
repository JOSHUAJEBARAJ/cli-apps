package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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

type results []result

func parseFlags() config {
	o := flag.String("o", "", "To send the output to the file")
	i := flag.Bool("i", false, "case sensitive")
	flag.Parse()
	if len(flag.Args()) < 1 {
		flag.PrintDefaults()
		os.Exit(1)
	}
	// if flag.NArg() == 0 {
	// 	flag.Usage()
	// 	os.Exit(1)
	// }
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

func search(r io.Reader, c config) ([]string, error) {
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
		return nil, err
	}
	return lines, nil
}

func printresult(c config, r result) {
	if c.output == "" {
		for _, line := range r.lines {
			if r.printwithfilename {
				fmt.Printf("%s : %s\n", c.filename, line)
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
			_, err := writer.WriteString(c.filename + ":" + line + "\n")
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
func main() {
	c := parseFlags()
	r := result{}
	// check if the filename is directory
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
			r.printwithfilename = true
			defer f.Close() // Close the file when we're done with it
			r.filename = c.filename
			r.lines, err = search(f, c)
			if err != nil {
				handleError(err)
			}

			printresult(c, r)
		}
		return
	}
	f, err := openFile(c)
	if err != nil {
		handleError(err)
	}
	defer f.Close() // Close the file when we're done with it

	r.lines, err = search(f, c)
	if err != nil {
		handleError(err)
	}
	printresult(c, r)
}

func handleError(err error) {
	fmt.Println(os.Stderr, "Error:", err)
	os.Exit(1)
}
