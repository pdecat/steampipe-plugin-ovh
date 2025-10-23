package ovh

import (
	"context"
	"errors"
	"net/http"

	"github.com/ovh/go-ovh/ovh"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

func connect(ctx context.Context, d *plugin.QueryData) (*ovh.Client, error) {
	// get ovh client from cache
	cacheKey := "ovh"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(*ovh.Client), nil
	}

	applicationKey := ""
	applicationSecret := ""
	consumerKey := ""
	endpoint := ""

	ovhConfig := GetConfig(d.Connection)

	if ovhConfig.ApplicationKey != nil {
		applicationKey = *ovhConfig.ApplicationKey
	}
	if ovhConfig.ApplicationSecret != nil {
		applicationSecret = *ovhConfig.ApplicationSecret
	}
	if ovhConfig.ConsumerKey != nil {
		consumerKey = *ovhConfig.ConsumerKey
	}
	if ovhConfig.Endpoint != nil {
		endpoint = *ovhConfig.Endpoint
	}

	if applicationKey == "" {
		return nil, errors.New("'application_key' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
	}
	if applicationSecret == "" {
		return nil, errors.New("'application_secret' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
	}
	if consumerKey == "" {
		return nil, errors.New("'consumer_key' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
	}
	if endpoint == "" {
		return nil, errors.New("'endpoint' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
	}

	client, err := ovh.NewClient(
		endpoint,
		applicationKey,
		applicationSecret,
		consumerKey,
	)
	if err != nil {
		plugin.Logger(ctx).Error("ovh.connect", "client_creation_error", err)
		return nil, err
	}

	// Set our custom transport for HTTP request/response tracing
	httpClient := &http.Client{
		Transport: http.DefaultTransport,
	}
	httpClient.Transport = NewTransport("OVH", ctx, httpClient.Transport, IsDebugOrHigher(ctx))
	client.Client = httpClient

	// Save to cache
	d.ConnectionManager.Cache.Set(cacheKey, client)

	return client, nil
}
