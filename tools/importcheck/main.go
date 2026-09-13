// Command importcheck reports Go files whose imports are not grouped as standard
// library, third-party, and first-party packages, in that order.
package main

import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type importGroup int

const (
	stdlib importGroup = iota
	thirdParty
	firstParty
)

type importInfo struct {
	path  string
	line  int
	group importGroup
	start int
	end   int
}

type issue struct {
	path string
	line int
	text string
}

func main() {
	root := flag.String("root", ".", "project root containing go.mod")
	flag.Parse()

	module, err := modulePath(filepath.Join(*root, "go.mod"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "importcheck: %v\n", err)
		os.Exit(2)
	}

	issues, err := checkGoFiles(*root, module)
	if err != nil {
		fmt.Fprintf(os.Stderr, "importcheck: %v\n", err)
		os.Exit(2)
	}
	for _, found := range issues {
		fmt.Printf("%s:%d: %s\n", found.path, found.line, found.text)
	}
	if len(issues) != 0 {
		os.Exit(1)
	}
}

func modulePath(path string) (string, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	for _, line := range strings.Split(string(contents), "\n") {
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[0] == "module" {
			return fields[1], nil
		}
	}
	return "", fmt.Errorf("%s does not contain a module directive", path)
}

func checkGoFiles(root, module string) ([]issue, error) {
	var issues []issue
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if path != root && (entry.Name() == ".git" || entry.Name() == "vendor" || strings.HasPrefix(entry.Name(), ".")) {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".go" {
			return nil
		}

		fileIssues, err := checkFile(path, module)
		if err != nil {
			return err
		}
		issues = append(issues, fileIssues...)
		return nil
	})
	return issues, err
}

func checkFile(path, module string) ([]issue, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, path, contents, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	imports := make([]importInfo, 0, len(file.Imports))
	for _, spec := range file.Imports {
		importPath, err := strconv.Unquote(spec.Path.Value)
		if err != nil {
			return nil, fmt.Errorf("%s: invalid import path: %w", path, err)
		}
		end := fileSet.Position(spec.End()).Offset
		imports = append(imports, importInfo{
			path:  importPath,
			line:  fileSet.Position(spec.Pos()).Line,
			group: classify(importPath, module),
			start: fileSet.Position(spec.Pos()).Offset,
			end:   end,
		})
	}

	if len(imports) == 0 {
		return nil, nil
	}

	var issues []issue
	expectedGroup := imports[0].group
	for index := 1; index < len(imports); index++ {
		previous, current := imports[index-1], imports[index]
		blankLine := strings.Contains(string(contents[previous.end:current.start]), "\n\n")
		if blankLine {
			if current.group == expectedGroup {
				issues = append(issues, issue{path, current.line, fmt.Sprintf("%q must be in the same %s group as the preceding import", current.path, groupName(current.group))})
			} else if current.group < expectedGroup {
				issues = append(issues, issue{path, current.line, fmt.Sprintf("%q is in the %s group after the %s group", current.path, groupName(current.group), groupName(expectedGroup))})
			}
			// A blank line starts a new group, so subsequent imports are
			// checked against this import's category rather than the previous
			// import's category.
			expectedGroup = current.group
			continue
		}

		if current.group > expectedGroup {
			issues = append(issues, issue{path, current.line, fmt.Sprintf("%q must be separated from the %s group by a blank line", current.path, groupName(expectedGroup))})
		} else if current.group < expectedGroup {
			issues = append(issues, issue{path, current.line, fmt.Sprintf("%q is in the %s group inside the %s group", current.path, groupName(current.group), groupName(expectedGroup))})
		}
	}
	return issues, nil
}

func classify(path, module string) importGroup {
	// Import paths without a domain are standard-library (and local relative
	// imports, if a package uses them). This deliberately does not use aliases:
	// the quoted import path is the authoritative package identity.
	firstPart := path
	if slash := strings.IndexByte(firstPart, '/'); slash >= 0 {
		firstPart = firstPart[:slash]
	}
	if !strings.Contains(firstPart, ".") {
		return stdlib
	}
	if path == module || strings.HasPrefix(path, module+"/") {
		return firstParty
	}
	return thirdParty
}

func groupName(group importGroup) string {
	switch group {
	case stdlib:
		return "stdlib"
	case thirdParty:
		return "third-party"
	default:
		return "first-party"
	}
}
