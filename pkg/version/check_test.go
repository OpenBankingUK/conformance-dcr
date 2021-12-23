package version

import (
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateCheck_OutdatedVersionUpdateAvailable(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte(`[{
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
	]`))
			if err != nil {
				t.Fail()
			}
		}))
	defer ts.Close()

	bb := NewGitHub(ts.URL)
	version := "0.0.2"

	upd, err := bb.UpdateAvailable(version)
	assert.NoError(t, err)
	assert.True(t, upd)
}

func TestUpdateCheck_NoTagsAvailable(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte(`[]`))
			if err != nil {
				t.Fail()
			}
		}))
	defer ts.Close()

	bb := NewGitHub(ts.URL)
	version := "0.0.2"

	upd, err := bb.UpdateAvailable(version)
	assert.NoError(t, err)
	assert.False(t, upd)
}

func TestUpdateCheck_UpToDateVersionNoUpdateAvailable(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte(`[{
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
	]`))
			if err != nil {
				t.Fail()
			}
		}))
	defer ts.Close()

	bb := NewGitHub(ts.URL)
	version := "1.3.0"

	upd, err := bb.UpdateAvailable(version)
	assert.NoError(t, err)
	assert.False(t, upd)
}

func TestUpdateCheck_NoLocalVersionSet(t *testing.T) {
	bb := NewGitHub("")
	version := ""

	update, err := bb.UpdateAvailable(version)
	expError := "version not set"
	assert.Equal(t, expError, err.Error())
	assert.False(t, update)
}

func TestUpdateCheck_LocalVersionInvalid(t *testing.T) {
	bb := NewGitHub("")
	version := "foobar"

	update, err := bb.UpdateAvailable(version)
	expError := "parse version: Malformed version: foobar"
	assert.Equal(t, expError, err.Error())
	assert.False(t, update)
}

func TestSortTags(t *testing.T) {
	actualTagList := tagList{
		{Name: "9.0.1-rc"},
		{Name: "0.0.4-dev"},
		{Name: "0.5.0"},
		{Name: "0.1.8"},
	}

	sort.Sort(actualTagList)

	expectedSorted := tagList([]tag{
		{Name: "0.0.4-dev"},
		{Name: "0.1.8"},
		{Name: "0.5.0"},
		{Name: "9.0.1-rc"},
	})

	assert.Equal(t, expectedSorted, actualTagList)
}
