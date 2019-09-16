package version

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestUpdateCheck_OutdatedVersionUpdateAvailable(t *testing.T) {}

func TestUpdateCheck_UpToDateVersionNoUpdateAvailable(t *testing.T) {}

func TestUpdateCheck_NoLocalVersionSet(t *testing.T) {
	bb := BitBucket{}
	version = ""

	msg, update, err := bb.UpdateCheck()
	expMessage := "Version check is unavailable at this time."
	assert.Equal(t, expMessage, msg)
	assert.Equal(t, false, update)
	expError := "no version found"
	assert.Equal(t, expError, err.Error())
}

func TestUpdateCheck_LocalVersionInvalid(t *testing.T) {}

func TestSortTags(t *testing.T) {
	actualTagList := tagList{
		{Name: "9.0.1-rc"},
		{Name: "0.0.4-dev"},
		{Name: "0.5.0"},
		{Name: "0.1.8"},
	}

	sort.Sort(actualTagList)

	expectedSorted := tagList([]Tag{
		{Name: "0.0.4-dev"},
		{Name: "0.1.8"},
		{Name: "0.5.0"},
		{Name: "9.0.1-rc"},
	})

	assert.Equal(t, expectedSorted, actualTagList)
}
