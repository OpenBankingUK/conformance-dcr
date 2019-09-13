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

// Checker defines functionality to reason about the current version of the software and if updates are available
type Checker interface {
	GetHumanVersion() string
	VersionFormatter(version string) (string, error)
	UpdateWarningVersion(version string) (string, bool, error)
}

// BitBucket helper with capability to get release versions from source control repository
type BitBucket struct {
	// bitBucketAPIRepository full URL of the TAG API 2.0 for the Conformance Suite.
	bitBucketAPIRepository string
}

// NewBitBucket returns a new instance of Checker.
func NewBitBucket(bitBucketAPIRepository string) BitBucket {
	return BitBucket{
		bitBucketAPIRepository: bitBucketAPIRepository,
	}
}

// Tag structure used map response of tags.
type Tag struct {
	Name          string `json:"name"`
	Date          string `json:"date"`
	CommitMessage string `json:"message"`
}

// TagsAPIResponse structure to map response.
type TagsAPIResponse struct {
	TagList []Tag `json:"values"`
}

func (t Tag) LessThan(subject string) bool {
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

type tagList []Tag

func (t tagList) Len() int {
	return len(t)
}

func (t tagList) Less(i, j int) bool {
	return t[i].LessThan(t[j].Name)
}

func (t tagList) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func getTags(body []byte) (*TagsAPIResponse, error) {
	var s = new(TagsAPIResponse)
	err := json.Unmarshal(body, &s)
	return s, err
}

// UpdateCheck checks the current version against the
// latest tag version on Bitbucket, if a newer version is found it
// returns a message and bool value that can be used to inform a user
// a newer version is available for download.
func (v BitBucket) UpdateCheck() (string, bool, error) {
	// A default message that can be presented to an end user.
	errorMessageUI := "Version check is unavailable at this time."

	// Some basic validation, check we have a version,
	if len(version) == 0 {
		return errorMessageUI, false, fmt.Errorf("no version found")
	}

	client := http.Client{
		Timeout: time.Duration(time.Second * 30),
	}

	// Try to get the latest tag using the BitBucket API.
	resp, err := client.Get(v.bitBucketAPIRepository)
	if err != nil {
		// If network error then return message, flag to NOT update and actual error.
		return errorMessageUI, false, errors.Wrap(err, "HTTP on GET to BitBucket API")
	}

	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errorMessageUI, false, errors.Wrap(err, "cannot read body API error.")
		}

		err = resp.Body.Close()
		if err != nil {
			return errorMessageUI, false, errors.Wrap(err, "close http response body")
		}

		s, err := getTags(body)
		if err != nil {
			return errorMessageUI, false, errors.Wrap(err, "getTags")
		}

		if len(s.TagList) == 0 {
			return errorMessageUI, false, errors.New("no Tags found")
		}

		// Convert the list of tags to tagList and sort
		tags := convertSortTags(s)

		// Get latest tag
		latestTag := tags[len(tags)-1].Name

		// Format version string to compare.
		versionLocal, err := goversion.NewVersion(version)
		versionRemote, err := goversion.NewVersion(latestTag)

		if versionLocal.LessThan(versionRemote) {
			errorMessageUI = fmt.Sprintf("Version %s of the this tool is out of date, please update to %s", versionLocal, versionRemote)
			return errorMessageUI, true, nil
		}
		// If local and remote version match or is higher then return false update flag.
		if versionLocal.GreaterThanOrEqual(versionRemote) {
			errorMessageUI = fmt.Sprintf("This tool is running the latest version %s", versionLocal)
			return errorMessageUI, false, nil
		}
	} else {
		// handle anything else other than 200 OK.
		return "", false, fmt.Errorf("HTTP status %d received", resp.StatusCode)
	}

	return errorMessageUI, false, nil
}

func convertSortTags(tar *TagsAPIResponse) tagList {
	tags := tagList{}
	for _, v := range tar.TagList {
		tags = append(tags, v)
	}
	sort.Sort(tags)
	return tags
}
