package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"unicode/utf8"

	"github.com/boyter/gocodewalker"
)

// injected when build
// go build -ldflags "-X 'main.version=$(git describe --tags 2>/dev/null || echo v0.0.0)' -X 'main.commit=$(git rev-parse --short HEAD 2>/dev/null || echo 0000000)'"
var (
	version = "dev"
	commit  = "none"
)

func isTextFile(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	// Read first 8000 bytes (like file(1) does)
	buf := make([]byte, 8000)
	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return false, err
	}
	buf = buf[:n]

	// If not valid UTF-8, treat as binary
	if !utf8.Valid(buf) {
		return false, nil
	}

	// Heuristic: if >10% of bytes are 0 or non-printable (except \n\r\t), treat as binary
	nonPrintable := 0
	for _, b := range buf {
		if (b < 32 && b != 9 && b != 10 && b != 13) || b == 0x7f {
			nonPrintable++
		}
	}
	if len(buf) > 0 && float64(nonPrintable)/float64(len(buf)) > 0.1 {
		return false, nil
	}
	return true, nil
}

func main() {
	listOnly := flag.Bool("list", false, "List matching files only (don't display content)")
	showVersion := flag.Bool("version", false, "Show version information")
	showHelp := flag.Bool("help", false, "Show help information")
	filterPattern := flag.String("filter", "", "Regex pattern to whitelist files (e.g. '\\.(go|java)$')")
	flag.Parse()

	if *showHelp {
		fmt.Println("Usage: catall [options]")
		fmt.Println("Concatenate and display contents of text files in directory tree respecting ignore files")
		fmt.Println("\nOptions:")
		flag.PrintDefaults()
		fmt.Println("\nExamples:")
		fmt.Println("  catall --filter '\\.(go|java)$'   # Process Go and Java files")
		fmt.Println("  catall --list --filter '\\.go$'    # List Go files")
		os.Exit(0)
	}

	if *showVersion {
		fmt.Printf("catall version %s (commit %s)\n", version, commit)
		os.Exit(0)
	}

	var filterRe *regexp.Regexp
	if *filterPattern != "" {
		var err error
		filterRe, err = regexp.Compile(*filterPattern)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error compiling regex: %v\n", err)
			fmt.Fprintf(os.Stderr, "Valid examples: '\\.go$', '\\.(java|kt)$', '^src/.*\\.java$'\n")
			os.Exit(1)
		}
	}

	// 1. Get output file name if redirected
	var outFile string
	if fi, err := os.Stdout.Stat(); err == nil && (fi.Mode()&os.ModeCharDevice) == 0 {
		// Output is redirected, try to get the file name from /proc/self/fd/1 (Linux)
		if link, err := os.Readlink("/proc/self/fd/1"); err == nil {
			outFile, _ = filepath.Abs(link)
		}
	}

	// 2. Walk files with gocodewalker
	fileListQueue := make(chan *gocodewalker.File, 100)
	fileWalker := gocodewalker.NewFileWalker(".", fileListQueue)
	fileWalker.IncludeHidden = true
	fileWalker.ExcludeDirectory = []string{".git", ".svn", ".hg"}
	// No extension filter, we want all files
	fileWalker.SetErrorHandler(func(e error) bool {
		// Print error and continue
		fmt.Fprintln(os.Stderr, "ERR", e.Error())
		return true
	})
	go fileWalker.Start()

	var files []string
	pwd, _ := os.Getwd()
	for f := range fileListQueue {
		absPath, err := filepath.Abs(f.Location)
		if err != nil {
			continue
		}
		// Skip output file itself
		if outFile != "" && absPath == outFile {
			continue
		}

		// exclude .gitignore and .ignore
		base := filepath.Base(absPath)
		if base == ".gitignore" || base == ".ignore" {
			continue
		}

		if filterRe != nil && !filterRe.MatchString(absPath) {
			continue
		}

		files = append(files, absPath)
	}

	// 3. Sort files alphabetically
	sort.Strings(files)

	if *listOnly {
		for _, absPath := range files {
			rel, _ := filepath.Rel(pwd, absPath)
			fmt.Println("./" + filepath.ToSlash(rel))
		}
		return
	}

	fmt.Println("Files not ignored in this directory:")
	for i, f := range files {
		rel, _ := filepath.Rel(pwd, f)
		fmt.Printf("%d: %s\n", i+1, "./"+filepath.ToSlash(rel))
	}

	// 4. Cat files one by one
	for _, absPath := range files {
		relPath, err := filepath.Rel(pwd, absPath)
		if err != nil {
			relPath = absPath
		}
		relPath = "./" + filepath.ToSlash(relPath)
		fmt.Printf("\n\n==== %s ====\n", relPath)

		isText, err := isTextFile(absPath)
		if err != nil || !isText {
			fmt.Println("BINARY OR BAD FORMAT")
			continue
		}

		// Print file content
		f, err := os.Open(absPath)
		if err != nil {
			fmt.Println("BINARY OR BAD FORMAT")
			continue
		}
		defer f.Close()

		// Use buffered reader for efficiency
		reader := bufio.NewReader(f)
		buf := make([]byte, 4096)
		for {
			n, err := reader.Read(buf)
			if n > 0 {
				os.Stdout.Write(buf[:n])
			}
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println("\nBINARY OR BAD FORMAT")
				break
			}
		}
	}
}
