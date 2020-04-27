package stencil

import (
	"bufio"
	"regexp"
	"strings"
)

// Markdown implements the markdown functionality.
type Markdown struct {
	*Stencil
}

// CopyMarkdownSnippets first treats the url as a markdown stripping out
// everything but code fences with names matching the provided regex.
// The name of the code fence is the first word after the triple backquote.
//
// Note that standard go templates are still executed on top of this
// so the embedded code can use any stencil function.
func (m *Markdown) CopyMarkdownSnippets(key, localPath, url, regex string) error {
	m.Printf("copying %s (snippets %s) to %s, key (%s)\n", url, regex, localPath, key)

	key = key + "(regex: " + regex + ")"
	m.Objects.addFile(key, localPath, url)

	data, err := m.executeFilter(url, func(md string) (string, error) {
		return m.FilterMarkdown(md, regex)
	})
	if err != nil {
		return m.Errorf("Error reading %s %v\n", url, err)
	}

	return m.Write(localPath, []byte(data), 0666)
}

func (m *Markdown) FilterMarkdown(data, regex string) (string, error) {
	result := ""
	if regex == "" {
		regex = ".*"
	}

	re, err := regexp.Compile(regex)
	if err != nil {
		return "", err
	}

	m.forEachCodeFence(data, func(name, code string) {
		if re.MatchString(name) {
			result += code
		}
	})

	return result, nil
}

func (m *Markdown) forEachCodeFence(md string, visit func(name, code string)) {
	re := regexp.MustCompile("(?msU)^```(.*)^```$")
	for _, fence := range re.FindAllString(md, -1) {
		fence = strings.Trim(fence, "`")
		scanner := bufio.NewScanner(strings.NewReader(fence))
		scanner.Scan()
		name := scanner.Text()
		visit(strings.TrimSpace(name), fence[len(name)+1:])
	}
}
