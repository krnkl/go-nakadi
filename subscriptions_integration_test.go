// +build integration

package nakadi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationSubscriptionAPI_Get(t *testing.T) {
	eventType := &EventType{}
	helperLoadTestData(t, "event-type-create.json", eventType)
	subscriptions := helperCreateSubscriptions(t, eventType, &Subscription{OwningApplication: "test-app", EventTypes: []string{eventType.Name}})
	defer helperDeleteSubscriptions(t, eventType, subscriptions...)

	client := New(defaultNakadiURL, &ClientOptions{ConnectionTimeout: time.Second})
	subAPI := NewSubscriptionAPI(client, &SubscriptionOptions{Retry: true})

	t.Run("fail not found", func(t *testing.T) {
		_, err := subAPI.Get("does-not-exist")

		require.Error(t, err)
		assert.Regexp(t, "does not exist", err)
	})

	t.Run("success", func(t *testing.T) {
		subscription, err := subAPI.Get(subscriptions[0].ID)

		require.NoError(t, err)
		assert.Equal(t, subscriptions[0].EventTypes, subscription.EventTypes)
		assert.Equal(t, subscriptions[0].OwningApplication, subscription.OwningApplication)
		assert.Equal(t, "default", subscription.ConsumerGroup)
	})
}

func TestIntegrationSubscriptionAPI_List(t *testing.T) {
	eventType := &EventType{}
	helperLoadTestData(t, "event-type-create.json", eventType)
	subscriptions := []*Subscription{
		{OwningApplication: "test-app", EventTypes: []string{eventType.Name}},
		{OwningApplication: "test-app2", EventTypes: []string{eventType.Name}}}
	subscriptions = helperCreateSubscriptions(t, eventType, subscriptions...)
	defer helperDeleteSubscriptions(t, eventType, subscriptions...)

	client := New(defaultNakadiURL, &ClientOptions{ConnectionTimeout: time.Second})
	subAPI := NewSubscriptionAPI(client, nil)

	fetched, err := subAPI.List()
	require.NoError(t, err)
	assert.Len(t, fetched, 2)
}

func TestIntegrationSubscriptionAPI_Create(t *testing.T) {
	eventType := &EventType{}
	helperLoadTestData(t, "event-type-create.json", eventType)
	helperCreateEventTypes(t, eventType)

	client := New(defaultNakadiURL, &ClientOptions{ConnectionTimeout: time.Second})
	subAPI := NewSubscriptionAPI(client, nil)

	t.Run("fail invalid subscription", func(t *testing.T) {
		_, err := subAPI.Create(&Subscription{OwningApplication: "test-api"})
		require.Error(t, err)
		assert.Regexp(t, "unable to create subscription", err)
	})

	t.Run("success", func(t *testing.T) {
		created, err := subAPI.Create(&Subscription{OwningApplication: "test-app", EventTypes: []string{eventType.Name}})
		require.NoError(t, err)
		assert.Equal(t, "test-app", created.OwningApplication)
		assert.Equal(t, []string{eventType.Name}, created.EventTypes)

		helperDeleteSubscriptions(t, eventType, created)
	})
}

func TestIntegrationSubscriptionAPI_Delete(t *testing.T) {
	eventType := &EventType{}
	helperLoadTestData(t, "event-type-create.json", eventType)
	subscriptions := helperCreateSubscriptions(t, eventType, &Subscription{OwningApplication: "test-app", EventTypes: []string{eventType.Name}})
	defer helperDeleteSubscriptions(t, eventType, subscriptions...)

	client := New(defaultNakadiURL, &ClientOptions{ConnectionTimeout: time.Second})
	subAPI := NewSubscriptionAPI(client, nil)

	t.Run("fail not found", func(t *testing.T) {
		err := subAPI.Delete("does-not-exist")

		require.Error(t, err)
		assert.Regexp(t, "does not exist", err)
	})

	t.Run("success", func(t *testing.T) {
		err := subAPI.Delete(subscriptions[0].ID)

		require.NoError(t, err)
	})
}

func helperCreateSubscriptions(t *testing.T, eventType *EventType, subscription ...*Subscription) []*Subscription {
	helperCreateEventTypes(t, eventType)
	var subscriptions []*Subscription
	for _, eventType := range subscription {
		serialized, err := json.Marshal(eventType)
		require.NoError(t, err)
		response, err := http.DefaultClient.Post(defaultNakadiURL+"/subscriptions", "application/json", bytes.NewReader(serialized))
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, response.StatusCode)
		created := &Subscription{}
		err = json.NewDecoder(response.Body).Decode(created)
		require.NoError(t, err)
		subscriptions = append(subscriptions, created)
	}
	return subscriptions
}

func helperDeleteSubscriptions(t *testing.T, eventType *EventType, subscriptions ...*Subscription) {
	for _, sub := range subscriptions {
		request, err := http.NewRequest("DELETE", defaultNakadiURL+"/subscriptions/"+sub.ID, nil)
		require.NoError(t, err)
		_, err = http.DefaultClient.Do(request)
		require.NoError(t, err)
	}
	helperDeleteEventTypes(t, eventType)
}
