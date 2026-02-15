package stringutil

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/google/uuid"
)

var tmpl = template.Must(template.ParseGlob("templates/*.html"))

func ParseStringToUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

func ParseStringsToUUIDs(ss []string) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	for _, s := range ss {
		id, err := uuid.Parse(s)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func ExtractUsernameFromEmail(email string) string {
	return email[:strings.Index(email, "@")]
}

func ConvertHtmlToString(name string, data any) (string, error) {
	var buf bytes.Buffer
	err := tmpl.ExecuteTemplate(&buf, name, data)
	return buf.String(), err
}
