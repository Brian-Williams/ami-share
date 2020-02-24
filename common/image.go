package common

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/rebuy-de/aws-nuke/pkg/types"
	"sort"
	"time"
)

type Image interface {
	Properties() types.Properties
	Date() time.Time
	String() string

	Match(Filter) bool
	AddTags(map[string]string, bool) error
	ShareWithAccount(string, bool) error
	CopyTags(*session.Session, bool) error
	MarshalYAML() (interface{}, error)
}

// For sorting images
// images must be sortable (by date) for determining most recent
type Images []Image

func (e Images) Len() int {
	return len(e)
}
func (e Images) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}
func (e Images) Less(i, j int) bool {
	return e[i].Date().Before(e[j].Date())
}

// Given a list of images apply a set of filters and pick the latest image
func ApplyFilters(images Images, filters []Filter) Images {
	var result Images
	for _, image := range images {
		matches := true
		for _, filter := range filters {
			matches = matches && image.Match(filter)
		}
		if matches {
			result = append(result, image)
		}
	}
	sort.Sort(result)
	if len(result) > 0 {
		// return the latest AMI only
		return Images{result[len(result)-1]}
	}
	return result
}
