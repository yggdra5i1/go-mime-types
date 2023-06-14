package mimego

import (
	"strings"

	"github.com/yggdra5i1/mimego/db"
	"github.com/yggdra5i1/mimego/utils"
)

type Mime struct {
	types      map[string]string
	extensions map[string]string
}

func (m *Mime) Define(typesMap map[string][]string, force bool) {
	for mimeType, extensions := range typesMap {
		extensions = utils.Map(extensions, strings.ToLower)

		mimeType = strings.ToLower(mimeType)

		for i := 0; i < len(extensions); i++ {
			ext := extensions[i]

			if string(ext[0]) == "*" {
				continue
			}

			m.types[ext] = mimeType
		}

		if _, ok := m.extensions[mimeType]; force || !ok {
			ext := extensions[0]
			if string(ext[0]) != "*" {
				m.extensions[mimeType] = ext
			} else {
				m.extensions[mimeType] = ext[1:]
			}
		}
	}
}

func buildMime(types []map[string][]string) *Mime {
	m := &Mime{}
	m.types = make(map[string]string)
	m.extensions = make(map[string]string)

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

func NewLite() *Mime {
	var types = []map[string][]string{db.StandardTypes}

	return buildMime(types)
}
