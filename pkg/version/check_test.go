package version

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"
)

func TestUpdateCheck_OutdatedVersionUpdateAvailable(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte(`{
	"values": [{
			"name": "0.0.1"
		},
		{
			"name": "0.0.0-dev"
		},
		{
			"name": "0.1.2-rc"
		},
		{
			"name": "0.0.4"
		},
		{
			"name": "1.3.0"
		},
		{
			"name": "0.0.2"
		}
	]
}`))
			if err != nil {
				t.Fail()
			}
		}))
	defer ts.Close()

	bb := BitBucket{
		bitBucketAPIRepository: ts.URL,
	}
	version = "0.0.2"

	msg, upd, err := bb.UpdateCheck()
	assert.NoError(t, err)
	assert.True(t, upd)
	expMsg := "Version 0.0.2 of the this tool is out of date, please update to 1.3.0"
	assert.Equal(t, expMsg, msg)
}

func TestUpdateCheck_UpToDateVersionNoUpdateAvailable(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte(`{
	"values": [{
			"name": "0.0.1"
		},
		{
			"name": "0.0.0-dev"
		},
		{
			"name": "0.1.2-rc"
		},
		{
			"name": "0.0.4"
		},
		{
			"name": "1.3.0"
		},
		{
			"name": "0.0.2"
		}
	]
}`))
			if err != nil {
				t.Fail()
			}
		}))
	defer ts.Close()

	bb := BitBucket{
		bitBucketAPIRepository: ts.URL,
	}
	version = "1.3.0"

	msg, upd, err := bb.UpdateCheck()
	assert.NoError(t, err)
	assert.False(t, upd)
	expMsg := "This tool is running the latest version 1.3.0"
	assert.Equal(t, expMsg, msg)
}

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
