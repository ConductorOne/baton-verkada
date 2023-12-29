package verkada

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const BaseUrl = "https://api.verkada.com"

type Client struct {
	httpClient *http.Client
	apiKey     string
}

type RequestBody struct {
	UserID string `json:"user_id"`
}

func NewClient(httpClient *http.Client, apiKey string) *Client {
	return &Client{
		httpClient: httpClient,
		apiKey:     apiKey,
	}
}

// ListUsers returns a list of all access users.
func (c *Client) ListUsers(ctx context.Context) ([]User, error) {
	url, _ := url.JoinPath(BaseUrl, "/access/v1/access_users")

	var res struct {
		Users []User `json:"access_members"`
	}

	if err := c.doRequest(ctx, http.MethodGet, url, &res, nil, nil); err != nil {
		return nil, err
	}

	return res.Users, nil
}

// GetUserAccessInformation returns user access information object.
func (c *Client) GetUserAccessInformation(ctx context.Context, userId string) (UserAccess, error) {
	accessUrl, _ := url.JoinPath(BaseUrl, "/access/v1/access_users/user")
	var res UserAccess

	q := url.Values{}
	q.Add("user_id", userId)

	if err := c.doRequest(ctx, http.MethodGet, accessUrl, &res, q, nil); err != nil {
		return UserAccess{}, err
	}

	return res, nil
}

// ListAccessGroups returns a list of all access groups.
func (c *Client) ListAccessGroups(ctx context.Context) ([]Group, error) {
	url, _ := url.JoinPath(BaseUrl, "/access/v1/access_groups")

	var res struct {
		Groups []Group `json:"access_groups"`
	}

	if err := c.doRequest(ctx, http.MethodGet, url, &res, nil, nil); err != nil {
		return nil, err
	}

	return res.Groups, nil
}

// AddUserToGroup adds user to access group.
func (c *Client) AddUserToGroup(ctx context.Context, groupId, userId string) error {
	groupUrl, _ := url.JoinPath(BaseUrl, "/access/v1/access_groups/group/user")

	var res struct {
		GroupID          string   `json:"group_id"`
		Name             string   `json:"name"`
		SuccessfulAdds   []string `json:"successful_adds"`
		UnsuccessfulAdds []string `json:"unsuccessful_adds"`
	}

	q := url.Values{}
	q.Add("group_id", groupId)

	requestBody := RequestBody{
		UserID: userId,
	}

	payload, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	if err := c.doRequest(ctx, http.MethodPut, groupUrl, &res, q, payload); err != nil {
		return err
	}

	if arrayContains(userId, res.UnsuccessfulAdds) {
		return fmt.Errorf("failed to add user to group")
	}

	return nil
}

// RemoveUserFromGroup removes user from access group.
func (c *Client) RemoveUserFromGroup(ctx context.Context, groupId, userId string) error {
	groupUrl, _ := url.JoinPath(BaseUrl, "/access/v1/access_groups/group/user")

	q := url.Values{}
	q.Add("group_id", groupId)
	q.Add("user_id", userId)

	if err := c.doRequest(ctx, http.MethodDelete, groupUrl, nil, q, nil); err != nil {
		return err
	}

	return nil
}

func (c *Client) doRequest(ctx context.Context, method, path string, res interface{}, params url.Values, payload []byte) error {
	req, err := http.NewRequestWithContext(ctx, method, path, bytes.NewReader(payload))
	if err != nil {
		return err
	}

	if params != nil {
		req.URL.RawQuery = params.Encode()
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("x-api-key", c.apiKey)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed with status: %s", resp.Status)
	}

	if method != http.MethodDelete {
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			return err
		}
	}

	return nil
}

func arrayContains(target string, array []string) bool {
	for _, item := range array {
		if target == item {
			return true
		}
	}
	return false
}
