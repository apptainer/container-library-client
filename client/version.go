// Copyright (c) 2019, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the LICENSE.md file
// distributed with the sources of this project regarding your rights to use or distribute this
// software.

package client

import (
	"context"
	"net/http"

	"github.com/blang/semver/v4"
	jsonresp "github.com/sylabs/json-resp"
)

const (
	// APIVersionV2Upload supports extended image upload functionality.
	APIVersionV2Upload = "2.0.0-alpha.1"
	// APIVersionV2ArchTags supports extended arch tags functionality.
	APIVersionV2ArchTags = "2.0.0-alpha.2"
)

// VersionInfo contains version information.
type VersionInfo struct {
	Version    string `json:"version"`
	APIVersion string `json:"apiVersion"`
}

// GetVersion gets version information from the Cloud-Library Service. The context controls the lifetime of
// the request.
func (c *Client) GetVersion(ctx context.Context) (vi VersionInfo, err error) {
	req, err := c.newRequest(ctx, http.MethodGet, "version", "", nil)
	if err != nil {
		return VersionInfo{}, err
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return VersionInfo{}, err
	}
	defer res.Body.Close()

	if err := jsonresp.ReadResponse(res.Body, &vi); err != nil {
		return VersionInfo{}, err
	}
	return vi, nil
}

// apiAtLeast returns true if cloud-library server supports requested (or greater) API version
func (c *Client) apiAtLeast(ctx context.Context, reqVersion string) bool {
	// query cloud-library server for supported api version
	vi, err := c.GetVersion(ctx)
	if err != nil || vi.APIVersion == "" {
		// unable to get cloud-library server API version, fallback to lowest
		// common denominator
		c.Logger.Logf("Unable to determine remote API version: %v", err)
		return false
	}
	v, err := semver.Make(vi.APIVersion)
	if err != nil {
		c.Logger.Logf("Unable to decode remote API version: %v", err)
		return false
	}
	minRequiredVers, err := semver.Make(reqVersion)
	if err != nil {
		c.Logger.Logf("Unable to decode minimum required version: %v", err)
		return false
	}
	return v.GTE(minRequiredVers)
}
