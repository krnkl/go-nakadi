package nakadi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/jarcoal/httpmock.v1"
)

func TestSubscription_Marshal(t *testing.T) {
	subscription := &Subscription{}
	expected := helperLoadTestData(t, "subscription.json", subscription)

	serialized, err := json.Marshal(subscription)
	require.NoError(t, err)
	require.NoError(t, err)
	assert.JSONEq(t, string(expected), string(serialized))
}

func TestSubscriptionAPI_Get(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	expected := &Subscription{}
	serialized := helperLoadTestData(t, "subscription.json", expected)

	client := &Client{nakadiURL: defaultNakadiURL, httpClient: http.DefaultClient}
	api := NewSubscriptionAPI(client, nil)
	url := fmt.Sprintf("%s/subscriptions/%s", defaultNakadiURL, expected.ID)

	t.Run("fail connection error", func(t *testing.T) {
		httpmock.RegisterResponder("GET", url, httpmock.NewErrorResponder(assert.AnError))

		_, err := api.Get(expected.ID)
		require.Error(t, err)
		assert.Regexp(t, assert.AnError, err)
	})

	t.Run("fail decode error", func(t *testing.T) {
		httpmock.RegisterResponder("GET", url, httpmock.NewStringResponder(http.StatusNotFound, ""))

		_, err := api.Get(expected.ID)
		require.Error(t, err)
		assert.Regexp(t, "unable to decode response body", err)
	})

	t.Run("fail with problem", func(t *testing.T) {
		problem := `{"detail": "not found"}`
		httpmock.RegisterResponder("GET", url, httpmock.NewStringResponder(http.StatusNotFound, problem))

		_, err := api.Get(expected.ID)
		require.Error(t, err)
		assert.Regexp(t, "not found", err)
	})

	t.Run("fail decode response", func(t *testing.T) {
		httpmock.RegisterResponder("GET", url, httpmock.NewStringResponder(http.StatusOK, ""))

		_, err := api.Get(expected.ID)
		require.Error(t, err)
		assert.Regexp(t, "unable to decode response body", err)
	})

	t.Run("success", func(t *testing.T) {
		httpmock.RegisterResponder("GET", url, httpmock.NewBytesResponder(http.StatusOK, serialized))

		requested, err := api.Get(expected.ID)
		require.NoError(t, err)
		assert.Equal(t, expected, requested)
	})
}

func TestSubscriptionAPI_List(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	expected := struct {
		Items []*Subscription `json:"items"`
	}{}
	helperLoadTestData(t, "subscriptions.json", &expected.Items)

	client := &Client{nakadiURL: defaultNakadiURL, httpClient: http.DefaultClient}
	api := NewSubscriptionAPI(client, nil)
	url := fmt.Sprintf("%s/subscriptions", defaultNakadiURL)

	t.Run("fail connection error", func(t *testing.T) {
		httpmock.RegisterResponder("GET", url, httpmock.NewErrorResponder(assert.AnError))

		_, err := api.List()
		require.Error(t, err)
		assert.Regexp(t, assert.AnError, err)
	})

	t.Run("fail decode error", func(t *testing.T) {
		httpmock.RegisterResponder("GET", url, httpmock.NewStringResponder(http.StatusInternalServerError, ""))

		_, err := api.List()
		require.Error(t, err)
		assert.Regexp(t, "unable to decode response body", err)
	})

	t.Run("fail with problem", func(t *testing.T) {
		problem := `{"detail": "not found"}`
		httpmock.RegisterResponder("GET", url, httpmock.NewStringResponder(http.StatusInternalServerError, problem))

		_, err := api.List()
		require.Error(t, err)
		assert.Regexp(t, "not found", err)
	})

	t.Run("fail decode response", func(t *testing.T) {
		httpmock.RegisterResponder("GET", url, httpmock.NewStringResponder(http.StatusOK, ""))

		_, err := api.List()
		require.Error(t, err)
		assert.Regexp(t, "unable to decode response body", err)
	})

	t.Run("success", func(t *testing.T) {
		responder, err := httpmock.NewJsonResponder(http.StatusOK, expected)
		require.NoError(t, err)
		httpmock.RegisterResponder("GET", url, responder)

		requested, err := api.List()
		require.NoError(t, err)
		assert.Equal(t, expected.Items, requested)
	})
}

func TestSubscriptionAPI_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	subscription := &Subscription{OwningApplication: "test-app", EventTypes: []string{"test-event.data"}}
	expected := &Subscription{}
	serialized := helperLoadTestData(t, "subscription.json", expected)

	client := &Client{nakadiURL: defaultNakadiURL, httpClient: http.DefaultClient}
	api := NewSubscriptionAPI(client, nil)
	url := fmt.Sprintf("%s/subscriptions", defaultNakadiURL)

	t.Run("fail connection error", func(t *testing.T) {
		httpmock.RegisterResponder("POST", url, httpmock.NewErrorResponder(assert.AnError))

		_, err := api.Create(subscription)
		require.Error(t, err)
		assert.Regexp(t, assert.AnError, err)
	})

	t.Run("fail with problem", func(t *testing.T) {
		problem := `{"detail": "not authorized"}`
		httpmock.RegisterResponder("POST", url, httpmock.NewStringResponder(http.StatusUnauthorized, problem))

		_, err := api.Create(subscription)
		require.Error(t, err)
		assert.Regexp(t, "not authorized", err)
	})

	t.Run("fail decode body", func(t *testing.T) {
		httpmock.RegisterResponder("POST", url, httpmock.NewStringResponder(http.StatusUnauthorized, ""))

		_, err := api.Create(subscription)
		require.Error(t, err)
		assert.Regexp(t, "unable to decode response body", err)

		httpmock.RegisterResponder("POST", url, httpmock.NewStringResponder(http.StatusOK, ""))

		_, err = api.Create(subscription)
		require.Error(t, err)
		assert.Regexp(t, "unable to decode response body", err)
	})

	t.Run("success", func(t *testing.T) {
		httpmock.RegisterResponder("POST", url, httpmock.Responder(func(r *http.Request) (*http.Response, error) {
			uploaded := &Subscription{}
			err := json.NewDecoder(r.Body).Decode(uploaded)
			require.NoError(t, err)
			assert.Equal(t, subscription, uploaded)
			return httpmock.NewBytesResponse(http.StatusOK, serialized), nil
		}))

		requested, err := api.Create(subscription)
		require.NoError(t, err)
		assert.Equal(t, expected, requested)
	})
}

func TestSubscriptionAPI_Delete(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	id := "7dd69d58-7f20-11e7-9748-133d6a0dbfb3"

	client := &Client{nakadiURL: defaultNakadiURL, httpClient: http.DefaultClient}
	api := NewSubscriptionAPI(client, nil)
	url := fmt.Sprintf("%s/subscriptions/%s", defaultNakadiURL, id)

	t.Run("fail connection error", func(t *testing.T) {
		httpmock.RegisterResponder("DELETE", url, httpmock.NewErrorResponder(assert.AnError))

		err := api.Delete(id)
		require.Error(t, err)
		assert.Regexp(t, assert.AnError, err)
	})

	t.Run("fail decode body", func(t *testing.T) {
		httpmock.RegisterResponder("DELETE", url, httpmock.NewStringResponder(http.StatusNotFound, ""))

		err := api.Delete(id)
		require.Error(t, err)
		assert.Regexp(t, "unable to decode response body", err)
	})

	t.Run("fail with problem", func(t *testing.T) {
		problem := `{"detail": "not found"}`
		httpmock.RegisterResponder("DELETE", url, httpmock.NewStringResponder(http.StatusNotFound, problem))

		err := api.Delete(id)
		require.Error(t, err)
		assert.Regexp(t, "not found", err)
	})

	t.Run("success", func(t *testing.T) {
		httpmock.RegisterResponder("DELETE", url, httpmock.NewStringResponder(http.StatusNoContent, ""))

		err := api.Delete(id)
		assert.NoError(t, err)
	})
}

func TestSubscriptionOptions_withDefaults(t *testing.T) {
	tests := []struct {
		Options  *SubscriptionOptions
		Expected *SubscriptionOptions
	}{
		{
			Options: nil,
			Expected: &SubscriptionOptions{
				InitialRetryInterval: defaultInitialRetryInterval,
				MaxRetryInterval:     defaultMaxRetryInterval,
				MaxElapsedTime:       defaultMaxElapsedTime,
			},
		},
		{
			Options: &SubscriptionOptions{InitialRetryInterval: time.Hour},
			Expected: &SubscriptionOptions{
				InitialRetryInterval: time.Hour,
				MaxRetryInterval:     defaultMaxRetryInterval,
				MaxElapsedTime:       defaultMaxElapsedTime,
			},
		},
		{
			Options: &SubscriptionOptions{MaxRetryInterval: time.Hour},
			Expected: &SubscriptionOptions{
				InitialRetryInterval: defaultInitialRetryInterval,
				MaxRetryInterval:     time.Hour,
				MaxElapsedTime:       defaultMaxElapsedTime,
			},
		},
		{
			Options: &SubscriptionOptions{MaxElapsedTime: time.Hour},
			Expected: &SubscriptionOptions{
				InitialRetryInterval: defaultInitialRetryInterval,
				MaxRetryInterval:     defaultMaxRetryInterval,
				MaxElapsedTime:       time.Hour,
			},
		},
		{
			Options: &SubscriptionOptions{Retry: true},
			Expected: &SubscriptionOptions{
				Retry:                true,
				InitialRetryInterval: defaultInitialRetryInterval,
				MaxRetryInterval:     defaultMaxRetryInterval,
				MaxElapsedTime:       defaultMaxElapsedTime,
			},
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.Expected, tt.Options.withDefaults())
	}
}
