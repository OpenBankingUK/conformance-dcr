package version

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"time"

	goversion "github.com/hashicorp/go-version"

	"github.com/pkg/errors"
)

type Checker interface {
	UpdateCheck() (bool, error)
}

// bitBucket helper with capability to get release versions from source control repository
type bitBucket struct {
	// bitBucketAPIRepository full URL of the TAG API 2.0 for the Conformance Suite.
	bitBucketAPIRepository string
}

// NewBitBucket returns a new instance of Checker.
func NewBitBucket(bitBucketAPIRepository string) Checker {
	return bitBucket{
		bitBucketAPIRepository: bitBucketAPIRepository,
	}
}

// tag structure used map response of tags.
type tag struct {
	Name          string `json:"name"`
	Date          string `json:"date"`
	CommitMessage string `json:"message"`
}

// tagsAPIResponse structure to map response.
type tagsAPIResponse struct {
	TagList []tag `json:"values"`
}

func (t tag) LessThan(subject string) bool {
	tv, err := goversion.NewVersion(t.Name)
	if err != nil {
		return false
	}
	sv, err := goversion.NewVersion(subject)
	if err != nil {
		return false
	}

	return tv.LessThan(sv)
}

type tagList []tag

func (t tagList) Len() int {
	return len(t)
}

func (t tagList) Less(i, j int) bool {
	return t[i].LessThan(t[j].Name)
}

func (t tagList) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func getTags(body []byte) (*tagsAPIResponse, error) {
	var s = new(tagsAPIResponse)
	err := json.Unmarshal(body, &s)
	return s, err
}

// UpdateCheck checks the current version against the
// latest tag version on Bitbucket, if a newer version is found true is returned,
// for all other cases, false is returned.
func (v bitBucket) UpdateCheck() (bool, error) {
	// Some basic validation, check we have a version,
	if len(version) == 0 {
		return false, errors.New("version not set")
	}

	// Format version string to compare.
	versionLocal, err := goversion.NewVersion(version)
	if err != nil {
		return false, errors.Wrap(err, "parse version")
	}

	tags, err := v.getTags()
	if err != nil {
		return false, errors.Wrap(err, "get tags from upstream repo")
	}

	// Convert the list of tags to tags and sort
	sort.Sort(tags)

	// Get latest tag
	latestTag := tags[len(tags)-1].Name

	versionRemote, err := goversion.NewVersion(latestTag)
	if err != nil {
		return false, errors.Wrap(err, "parse latest tag")
	}

	if versionLocal.LessThan(versionRemote) {
		return true, nil
	}
	// If local and remote version match or is higher then return false update flag.
	if versionLocal.GreaterThanOrEqual(versionRemote) {
		return false, nil
	}

	return false, nil
}

func (v *bitBucket) getTags() (tagList, error) {
	client := http.Client{
		Timeout: time.Second * 30,
	}

	// Try to get the latest tag using the BitBucket API.
	resp, err := client.Get(v.bitBucketAPIRepository)
	if err != nil {
		// If network error then return message, flag to NOT update and actual error.
		return nil, errors.Wrap(err, "HTTP on GET to BitBucket API")
	}

	if resp.StatusCode != http.StatusOK {
		// handle anything else other than 200 OK.
		return nil, fmt.Errorf("HTTP status %d received", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "cannot read body API error.")
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, errors.Wrap(err, "close http response body")
	}

	s, err := getTags(body)
	if err != nil {
		return nil, errors.Wrap(err, "getTags")
	}

	tags := tagList{}
	for _, v := range s.TagList {
		tags = append(tags, v)
	}

	return tags, nil
}
