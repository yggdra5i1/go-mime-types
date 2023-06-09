package mimego

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/yggdra5i1/mimego/db"
	"github.com/yggdra5i1/mimego/utils"
)

const (
	redefineErrMsg = `Attempt to change mapping for \"%s\" extension from \"%s\" to \"%s\". 
					  Pass force=true to allow this, otherwise remove \"%s\" from the list 
					  of extensions for \"%s\".`
)

const (
	textMediaType  = "text"
	imageMediaType = "image"
	audioMediaType = "audio"
	videoMediaType = "video"
)

var (
	defaultTypes = []map[string][]string{
		db.StandardTypes,
		db.OtherTypes,
	}
)

type Mime struct {
	types      map[string]string
	extensions map[string][]string
}

func (m *Mime) defineTypeForExtensions(mimeType string, extensions []string, force bool) {
	for _, ext := range extensions {
		if string(ext[0]) == "*" {
			continue
		}
		if _, ok := m.types[ext]; !force && ok {
			panic(fmt.Sprintf(redefineErrMsg, ext, m.types[ext], mimeType, ext, mimeType))
		}
		m.types[ext] = mimeType
	}
}

func (m *Mime) defineExtensionsForType(mimeType string, extensions []string) {
	if _, ok := m.extensions[mimeType]; !ok {
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

func (m *Mime) Define(typeMap map[string][]string, force bool) {
	for mimeType, extensions := range typeMap {
		extensions = utils.Map(extensions, strings.ToLower)
		mimeType = strings.ToLower(mimeType)

		m.defineTypeForExtensions(mimeType, extensions, force)
		m.defineExtensionsForType(mimeType, extensions)
	}
}

func buildMime(types []map[string][]string) *Mime {
	m := &Mime{
		types:      make(map[string]string),
		extensions: make(map[string][]string),
	}

	for _, t := range types {
		m.Define(t, false)
	}

	return m
}

func New(types []map[string][]string) *Mime {
	if len(types) > 0 {
		return buildMime(types)
	}

	return buildMime(defaultTypes)
}

func Lite(types []map[string][]string) *Mime {
	if len(types) > 0 {
		return buildMime(types)
	}

	return buildMime([]map[string][]string{db.StandardTypes})
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
	re := regexp.MustCompile(`^\s*([^;\s]*)`)
	typeMatch := re.FindStringSubmatch(mimeType)
	if len(typeMatch) > 1 {
		mimeType = typeMatch[1]
	}
	exts, ok := m.extensions[strings.ToLower(mimeType)]
	if ok {
		return exts
	}
	return nil
}

func IsText(mimeType string) bool {
	return getMediaType(mimeType) == textMediaType
}

func IsImage(mimeType string) bool {
	return getMediaType(mimeType) == imageMediaType
}

func IsAudio(mimType string) bool {
	return getMediaType(mimType) == audioMediaType
}

func IsVideo(mimeType string) bool {
	return getMediaType(mimeType) == videoMediaType
}

func getMediaType(mimeType string) string {
	return strings.Split(mimeType, "/")[0]
}
