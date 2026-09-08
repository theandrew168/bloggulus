// Command csscheck reports CSS classes that are not used by the HTML templates and
// classes used by templates that have no CSS definition.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type location struct {
	path string
	line int
}

type classLocations map[string][]location

var (
	cssClassRE      = regexp.MustCompile(`\.([-_a-zA-Z][-_a-zA-Z0-9]*)`)
	htmlClassAttrRE = regexp.MustCompile(`(?i)\bclass\s*=\s*["']([^"']*)["']`)
)

func main() {
	cssDir := flag.String("css", "public/css", "directory containing CSS files")
	templateDir := flag.String("templates", "backend/web", "directory containing HTML templates")
	flag.Parse()

	cssClasses, err := collectCSSClasses(*cssDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "csscheck: %v\n", err)
		os.Exit(2)
	}
	htmlClasses, err := collectHTMLClasses(*templateDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "csscheck: %v\n", err)
		os.Exit(2)
	}

	unused := difference(cssClasses, htmlClasses)
	missing := difference(htmlClasses, cssClasses)
	if len(unused) > 0 {
		printReport("Unused CSS classes", unused, cssClasses)
	}
	if len(missing) > 0 {
		printReport("Missing CSS classes", missing, htmlClasses)
	}

	if len(unused) != 0 || len(missing) != 0 {
		os.Exit(1)
	}
}

func collectCSSClasses(root string) (classLocations, error) {
	classes := make(classLocations)
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || strings.ToLower(filepath.Ext(path)) != ".css" {
			return nil
		}

		contents, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		for _, found := range cssClassMatches(stripCSSComments(string(contents))) {
			classes[found.name] = append(classes[found.name], location{path: path, line: lineNumber(string(contents), found.offset)})
		}
		return nil
	})
	return classes, err
}

type match struct {
	name   string
	offset int
}

// cssClassMatches only examines selector preludes (the text immediately before
// an opening brace), so strings such as ".example" in a declaration are ignored.
func cssClassMatches(contents string) []match {
	var matches []match
	start := 0
	for index, char := range contents {
		switch char {
		case '{':
			selector := contents[start:index]
			for _, found := range cssClassRE.FindAllStringSubmatchIndex(selector, -1) {
				matches = append(matches, match{name: selector[found[2]:found[3]], offset: start + found[0]})
			}
			start = index + 1
		case '}':
			start = index + 1
		}
	}
	return matches
}

func collectHTMLClasses(root string) (classLocations, error) {
	classes := make(classLocations)
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || strings.ToLower(filepath.Ext(path)) != ".html" {
			return nil
		}

		contents, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		text := string(contents)
		for _, attr := range htmlClassAttrRE.FindAllStringSubmatchIndex(text, -1) {
			valueStart, valueEnd := attr[2], attr[3]
			value := text[valueStart:valueEnd]
			for _, field := range strings.Fields(value) {
				offset := valueStart + strings.Index(value, field)
				classes[field] = append(classes[field], location{path: path, line: lineNumber(text, offset)})
			}
		}
		return nil
	})
	return classes, err
}

func stripCSSComments(contents string) string {
	commentRE := regexp.MustCompile(`(?s)/\*.*?\*/`)
	return commentRE.ReplaceAllStringFunc(contents, func(comment string) string {
		// Replace bytes rather than runes so offsets in the stripped string still
		// refer to the same locations in the original file.
		bytes := []byte(comment)
		for i, char := range bytes {
			if char != '\n' && char != '\r' {
				bytes[i] = ' '
			}
		}
		return string(bytes)
	})
}

func lineNumber(contents string, offset int) int {
	return 1 + strings.Count(contents[:offset], "\n")
}

func difference(left, right classLocations) []string {
	var result []string
	for name := range left {
		if _, exists := right[name]; !exists {
			result = append(result, name)
		}
	}
	sort.Strings(result)
	return result
}

func printReport(title string, names []string, locations classLocations) {
	fmt.Println(title + ":")
	for _, name := range names {
		refs := append([]location(nil), locations[name]...)
		sort.Slice(refs, func(i, j int) bool {
			if refs[i].path == refs[j].path {
				return refs[i].line < refs[j].line
			}
			return refs[i].path < refs[j].path
		})
		parts := make([]string, len(refs))
		for i, ref := range refs {
			parts[i] = fmt.Sprintf("%s:%d", ref.path, ref.line)
		}
		fmt.Printf("  .%s (%s)\n", name, strings.Join(parts, ", "))
	}
}
