package registryclient

import (
	"regexp"
	"sort"
	"strings"

	"github.com/hashicorp/go-version"
)

func sortTags(tags []string) {
	sort.Slice(tags, func(i, j int) bool {
		ti, tj := tags[i], tags[j]

		iHasNum := hasNumber(ti)
		jHasNum := hasNumber(tj)

		if !iHasNum && jHasNum {
			return true
		}

		if iHasNum && !jHasNum {
			return false
		}

		if !iHasNum && !jHasNum {
			if ti == "latest" {
				return true
			}

			if tj == "latest" {
				return false
			}

			return ti < tj
		}

		vStrI := strings.Split(ti, "-")[0]
		vStrJ := strings.Split(tj, "-")[0]

		vi, errI := version.NewVersion(vStrI)
		vj, errJ := version.NewVersion(vStrJ)

		if err1, err2 := errI == nil, errJ == nil; err1 && err2 {
			if !vi.Equal(vj) {
				return vi.GreaterThan(vj)
			}

			return strings.Count(ti, "-") < strings.Count(tj, "-")
		}

		return ti < tj
	})
}

var rgxHasNumber = regexp.MustCompile(`[0-9]`)

func hasNumber(s string) bool {
	return rgxHasNumber.MatchString(s)
}
