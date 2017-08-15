package gocd

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestConfiguration(t *testing.T) {
	setup()
	defer teardown()

	t.Run("HasAuth", testConfigurationHasAuth)
	t.Run("New", testConfigurationNew)
	t.Run("SanitizeURL", testConfigurationSantizieURL)
	t.Run("GetVersion", testConfigurationGetVersion)
}

func testConfigurationGetVersion(t *testing.T) {
	mux.HandleFunc("/api/version", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET", "Unexpected HTTP method")
		j, _ := ioutil.ReadFile("test/resources/version.0.json")
		fmt.Fprint(w, string(j))
	})
	v, _, err := client.Configuration.GetVersion(context.Background())
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "https://build.go.cd/go/api/version", v.Links.Self.String())
	assert.Equal(t, "https://api.gocd.org/#version", v.Links.Doc.String())
	assert.Equal(t, "16.6.0", v.Version)
	assert.Equal(t, "3348", v.BuildNumber)
	assert.Equal(t, "a7a5717cbd60c30006314fb8dd529796c93adaf0", v.GitSHA)
	assert.Equal(t, "16.6.0 (3348-a7a5717cbd60c30006314fb8dd529796c93adaf0)", v.FullVersion)
	assert.Equal(t, "https://github.com/gocd/gocd/commits/a7a5717cbd60c30006314fb8dd529796c93adaf0", v.CommitURL)
}

func testConfigurationSantizieURL(t *testing.T) {
	u := sanitizeURL(nil)
	assert.Nil(t, u)
}

func testConfigurationNew(t *testing.T) {
	c := Configuration{}
	client := c.Client()
	assert.NotNil(t, client)
}

func testConfigurationHasAuth(t *testing.T) {
	c := Configuration{}

	c.Username = "user"
	c.Password = "pass"
	assert.True(t, c.HasAuth())

	c.Username = "user"
	c.Password = ""
	assert.False(t, c.HasAuth())

	c.Username = ""
	c.Password = "pass"
	assert.False(t, c.HasAuth())

	c.Username = ""
	c.Password = ""
	assert.False(t, c.HasAuth())
}
