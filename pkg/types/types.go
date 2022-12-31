package types

import (
	"crypto/sha1"
	"encoding/base64"
	"path"
	"regexp"
	"strings"
)

//go:generate protoc --proto_path=../../proto --go_out=. --go_opt=paths=source_relative ../../proto/models.proto

type ClientsArr []*Client

func (arr ClientsArr) ById(id string) *Client {
	for _, a := range arr {
		if a == nil {
			continue
		}
		if a.Id == id {
			return a
		}
	}
	return nil
}

func (arr ClientsArr) ByName(name string) *Client {
	for _, a := range arr {
		if a == nil {
			continue
		}
		if a.Name == name {
			return a
		}
	}
	return nil
}

func (a *Client) MakeId(salt string) string {
	buf := sha1.Sum([]byte(a.Description + a.Name + salt))
	return base64.RawStdEncoding.EncodeToString(buf[:])
}

// Courtesy of https://github.com/kennygrant/sanitize
var (
	separators  = regexp.MustCompile(`[ &_=+:]`)
	dashes      = regexp.MustCompile(`[\/-]+`)
	illegalPath = regexp.MustCompile(`[^[:alnum:]\~\-\./]`)
)

func SanitizePath(s string) string {
	filePath := strings.ToLower(s)
	filePath = strings.Replace(filePath, "..", "", -1)
	filePath = path.Clean(filePath)
	filePath = strings.Trim(filePath, " ")
	filePath = separators.ReplaceAllString(filePath, "-")
	filePath = illegalPath.ReplaceAllString(filePath, "-")
	filePath = dashes.ReplaceAllString(filePath, "-")
	return filePath
}

func (a *Template) MakeId() string {
	return SanitizePath(a.What)
}
