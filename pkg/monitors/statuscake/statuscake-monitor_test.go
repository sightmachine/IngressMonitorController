package statuscake

import (
	"os"
	"testing"

	"github.com/stakater/IngressMonitorController/pkg/config"
	"github.com/stakater/IngressMonitorController/pkg/models"
	"github.com/stakater/IngressMonitorController/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestAddMonitorWithCorrectValues(t *testing.T) {
	config := config.GetControllerConfig()

	service := StatusCakeMonitorService{}
	provider := util.GetProviderWithName(config, "StatusCake")
	if provider == nil {
		panic("Failed to find provider")
	}
	service.Setup(*provider)

	m := models.Monitor{Name: "google-test", URL: "https://google1.com"}
	service.Add(m)

	mRes, err := service.GetByName("google-test")

	if err != nil {
		t.Error("Error: " + err.Error())
	}
	if mRes.Name != m.Name || mRes.URL != m.URL {
		t.Error("URL and name should be the same")
	}
	service.Remove(*mRes)

	monitor, err := service.GetByName(mRes.Name)

	if monitor != nil {
		t.Error("Monitor should've been deleted ", monitor, err)
	}
}

func TestUpdateMonitorWithCorrectValues(t *testing.T) {
	config := config.GetControllerConfig()

	service := StatusCakeMonitorService{}

	provider := util.GetProviderWithName(config, "StatusCake")
	if provider == nil {
		panic("Failed to find provider")
	}
	service.Setup(*provider)

	m := models.Monitor{Name: "google-test", URL: "https://google.com"}
	service.Add(m)

	mRes, err := service.GetByName("google-test")

	if err != nil {
		t.Error("Error: " + err.Error())
	}
	if mRes.Name != m.Name || mRes.URL != m.URL {
		t.Error("URL and name should be the same")
	}

	mRes.URL = "https://facebook.com"

	service.Update(*mRes)

	mRes, err = service.GetByName("google-test")

	if err != nil {
		t.Error("Error: " + err.Error())
	}
	if mRes.URL != "https://facebook.com" {
		t.Error("URL and name should be the same")
	}

	service.Remove(*mRes)

	monitor, err := service.GetByName(mRes.Name)

	if monitor != nil {
		t.Error("Monitor should've been deleted ", monitor, err)
	}
}

func TestBuildUpsertFormAnnotations(t *testing.T) {
	m := models.Monitor{Name: "google-test", URL: "https://google.com"}
	m.Annotations = map[string]string{
		"statuscake.monitor.stakater.com/check-rate":       "60",
		"statuscake.monitor.stakater.com/test-type":        "TCP",
		"statuscake.monitor.stakater.com/paused":           "true",
		"statuscake.monitor.stakater.com/ping-url":         "",
		"statuscake.monitor.stakater.com/follow-redirect":  "true",
		"statuscake.monitor.stakater.com/port":             "7070",
		"statuscake.monitor.stakater.com/trigger-rate":     "1",
		"statuscake.monitor.stakater.com/contact-group":    "123456,654321",
		"statuscake.monitor.stakater.com/basic-auth-user":  "testuser",
		"statuscake.monitor.stakater.com/node-locations":   "",
		"statuscake.monitor.stakater.com/status-codes":     "500,501,502,503,504,505",
		"statuscake.monitor.stakater.com/confirmation":     "2",
		"statuscake.monitor.stakater.com/enable-ssl-alert": "true",
		"statuscake.monitor.stakater.com/test-tags":        "test,testrun,uptime",
		"statuscake.monitor.stakater.com/real-browser":     "true",
	}

	oldEnv := os.Getenv("testuser")
	os.Setenv("testuser", "testpass")
	defer os.Setenv("testuser", oldEnv)

	vals := buildUpsertForm(m, "")
	assert.Equal(t, "testuser", vals.Get("BasicUser"))
	assert.Equal(t, "testpass", vals.Get("BasicPass"))
	assert.Equal(t, "60", vals.Get("CheckRate"))
	assert.Equal(t, "2", vals.Get("Confirmation"))
	assert.Equal(t, "123456,654321", vals.Get("ContactGroup"))
	assert.Equal(t, "1", vals.Get("EnableSSLAlert"))
	assert.Equal(t, "1", vals.Get("FollowRedirect"))
	assert.Equal(t, "", vals.Get("NodeLocations"))
	assert.Equal(t, "1", vals.Get("Paused"))
	assert.Equal(t, "", vals.Get("PingURL"))
	assert.Equal(t, "7070", vals.Get("Port"))
	assert.Equal(t, "1", vals.Get("RealBrowser"))
	assert.Equal(t, "500,501,502,503,504,505", vals.Get("StatusCodes"))
	assert.Equal(t, "test,testrun,uptime", vals.Get("TestTags"))
	assert.Equal(t, "TCP", vals.Get("TestType"))
	assert.Equal(t, "1", vals.Get("TriggerRate"))
}
