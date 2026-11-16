package ovh

import (
	"context"
	"net/http"

	"github.com/ovh/go-ovh/ovh"
)

// listAllV2 streams every item of an OVHcloud API v2 listing to stream, following the cursor the
// API returns in the X-Pagination-Cursor-Next response header. A plain client.Get only returns the
// first page of such a listing.
func listAllV2[T any](ctx context.Context, client *ovh.Client, path string, stream func(T)) error {
	cursor := ""
	for {
		req, err := client.NewRequest(http.MethodGet, path, nil, true)
		if err != nil {
			return err
		}
		req = req.WithContext(ctx)
		if cursor != "" {
			req.Header.Set("X-Pagination-Cursor", cursor)
		}

		resp, err := client.Do(req)
		if err != nil {
			return err
		}

		var page []T
		if err := client.UnmarshalResponse(resp, &page); err != nil {
			return err
		}
		for _, item := range page {
			stream(item)
		}

		cursor = resp.Header.Get("X-Pagination-Cursor-Next")
		if cursor == "" {
			return nil
		}
	}
}
