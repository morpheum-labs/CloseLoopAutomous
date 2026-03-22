package httpapi

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var describeCommitsRE = regexp.MustCompile(`^(.+)-(\d+)-g[0-9a-fA-F]+$`)
var semverCoreRE = regexp.MustCompile(`^v?(\d+(?:\.\d+){1,2})`)

type versionResponse struct {
	Version         string `json:"version"`
	Tag             string `json:"tag"`
	Number          string `json:"number"`
	CommitsAfterTag int    `json:"commits_after_tag"`
	Commit          string `json:"commit"`
	Dirty           bool   `json:"dirty"`
}

func buildVersionResponse(version, commit string) versionResponse {
	v := strings.TrimSpace(version)
	dirty := strings.HasSuffix(v, "-dirty")
	if dirty {
		v = strings.TrimSuffix(v, "-dirty")
	}
	tag := v
	commitsAfter := 0
	if m := describeCommitsRE.FindStringSubmatch(v); len(m) == 3 {
		tag = m[1]
		commitsAfter, _ = strconv.Atoi(m[2])
	}
	num := semverNumberFromTag(tag)
	return versionResponse{
		Version:         strings.TrimSpace(version),
		Tag:             tag,
		Number:          num,
		CommitsAfterTag: commitsAfter,
		Commit:          strings.TrimSpace(commit),
		Dirty:           dirty,
	}
}

func semverNumberFromTag(tag string) string {
	tag = strings.TrimSpace(tag)
	if m := semverCoreRE.FindStringSubmatch(tag); len(m) > 1 {
		return m[1]
	}
	return ""
}

func (h *Handlers) version(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, buildVersionResponse(h.BuildVersion, h.BuildCommit))
}
