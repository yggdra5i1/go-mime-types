package mimego

import (
	"regexp"
	"strings"

	"github.com/yggdra5i1/mimego/db"
	"github.com/yggdra5i1/mimego/utils"
)

type Mime struct {
	types      map[string]string
	extensions map[string][]string
}

func (m *Mime) defineTypeForExtensions(mimeType string, extensions []string) {
	for _, ext := range extensions {
		if string(ext[0]) == "*" {
			continue
		}

		m.types[ext] = mimeType
	}
}

func (m *Mime) defineExtensionsForType(mimeType string, extensions []string, force bool) {
	if _, ok := m.extensions[mimeType]; force || !ok {
		extensions = utils.Map(extensions, func(s string) string {
			if string(s[0]) != "*" {
				return s
			} else {
				return s[1:]
			}
		})

		m.extensions[mimeType] = extensions
	}
}

func (m *Mime) Define(typesMap map[string][]string, force bool) {
	for mimeType, extensions := range typesMap {
		extensions = utils.Map(extensions, strings.ToLower)
		mimeType = strings.ToLower(mimeType)

		m.defineTypeForExtensions(mimeType, extensions)
		m.defineExtensionsForType(mimeType, extensions, force)
	}
}

func buildMime(types []map[string][]string) *Mime {
	m := &Mime{}
	m.types = make(map[string]string)
	m.extensions = make(map[string][]string)

	for _, t := range types {
		m.Define(t, false)
	}

	return m
}

func New() *Mime {
	var types = []map[string][]string{
		db.StandardTypes,
		db.OtherTypes,
	}

	return buildMime(types)
}

func Lite() *Mime {
	var types = []map[string][]string{db.StandardTypes}

	return buildMime(types)
}

// Lookup a mime type based on extension
func (m *Mime) GetType(path string) (string, bool) {
	lastRegexp := regexp.MustCompile(`^.*[/\\]`)
	extRegexp := regexp.MustCompile(`^.*\.`)

	last := strings.ToLower(lastRegexp.ReplaceAllString(path, ""))
	ext := strings.ToLower(extRegexp.ReplaceAllString(last, ""))

	hasPath := len(last) < len(path)
	hasDot := len(ext) < len(last)-1

	if hasDot || !hasPath {
		mimeType, ok := m.types[ext]
		return mimeType, ok
	}

	return "", false
}

func (m *Mime) GetExtensions(mimeType string) []string {
	exts, _ := m.extensions[mimeType]
	return exts
}
