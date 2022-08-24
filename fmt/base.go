package fmt

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"gofimpt/errors"
)

type Import struct {
	comments []string
	alias    string
	pkg      string
}

// IsModule check if pkg starts with module name
func (i Import) IsModule(module string) bool {
	if len(module) == 0 {
		return false
	}
	return strings.HasPrefix(i.pkg, module)
}

// IsBuiltin check if pkg contains dot, note that you should call IsModule before calling this
func (i Import) IsBuiltin() bool {
	return !strings.Contains(i.pkg, ".")
}

type importSorter []Import

func (s importSorter) Len() int {
	return len(s)
}

func (s importSorter) Less(i, j int) bool {
	return s[i].pkg < s[j].pkg
}

func (s importSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (i Import) String() string {
	var result []string
	for _, comment := range i.comments {
		result = append(result, fmt.Sprintf("\t%s", comment))
	}
	imp := fmt.Sprintf("\t\"%s\"", i.pkg)
	if len(i.alias) != 0 {
		imp = fmt.Sprintf("\t%s \"%s\"", i.alias, i.pkg)
	}
	result = append(result, imp)
	return strings.Join(result, "\n")
}

func fmtImport(module, content string) string {
	var lines []string
	for _, s := range strings.Split(content, "\n") {
		s = strings.TrimSpace(s)
		s = strings.Trim(s, `"`)
		if len(s) == 0 {
			continue
		}
		lines = append(lines, s)
	}
	isComment := func(line string) bool {
		return strings.HasPrefix(strings.TrimSpace(line), "//")
	}
	var imports []Import
	imp := Import{}
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if isComment(line) {
			imp.comments = append(imp.comments, line)
			continue
		}
		ss := strings.Split(line, " ")
		switch {
		case len(ss) == 1:
			imp.pkg = line
		case len(ss) == 2:
			imp.alias = ss[0]
			imp.pkg = strings.Trim(ss[1], `"`)
		default:
			fmt.Printf("too many spaces in \"%s\"\n", line)
		}
		imports = append(imports, imp)
		imp = Import{}
	}
	var builtins []Import
	var thirdParties []Import
	var modules []Import
	for _, imp := range imports {
		switch {
		case imp.IsModule(module):
			modules = append(modules, imp)
		case imp.IsBuiltin():
			builtins = append(builtins, imp)
		default:
			thirdParties = append(thirdParties, imp)
		}
	}
	sort.Sort(importSorter(builtins))
	sort.Sort(importSorter(thirdParties))
	sort.Sort(importSorter(modules))
	lines = lines[:0]
	for _, builtin := range builtins {
		lines = append(lines, builtin.String())
	}
	if len(thirdParties) != 0 {
		lines = append(lines, "")
		for _, thirdParty := range thirdParties {
			lines = append(lines, thirdParty.String())
		}
	}
	if len(modules) != 0 {
		lines = append(lines, "")
		for _, module := range modules {
			lines = append(lines, module.String())
		}
	}
	return strings.Join(lines, "\n")
}

func FormatCode(module, data string) string {
	token := "import (\n"
	segments := strings.Split(data, token)
	// no import group
	if len(segments) == 1 {
		return data
	}
	result := []string{segments[0]}
	for _, seg := range segments[1:] {
		ss := strings.SplitN(seg, ")", 2)
		if len(ss) != 2 {
			result = append(result, seg)
			continue
		}
		imp, other := ss[0], ss[1]
		imp = fmtImport(module, imp)
		seg = fmt.Sprintf("%s\n)%s", imp, other)
		result = append(result, seg)
	}
	return strings.Join(result, token)
}

func FormatFile(module, file string) error {
	fmt.Println("processing", file)
	data, err := os.ReadFile(file)
	if err != nil {
		return errors.Trace(err)
	}
	result := FormatCode(module, string(data))
	f, _ := os.Stat(file)
	if result == string(data) {
		// nothing changed
		return nil
	}
	err = os.WriteFile(file, []byte(result), f.Mode())
	if err != nil {
		return errors.Trace(err)
	}
	fmt.Println("formatted", file)
	return nil
}
