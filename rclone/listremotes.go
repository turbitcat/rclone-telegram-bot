package rclone

import (
	"encoding/json"
	"fmt"
)

// json struct:
// operations/copyurl
type ListremotesResponse struct {
	Remotes []string
}

// ListRemotes returns a list of avalible remotes.
func (rs *RcloneServer) ListRemotes() ([]string, error) {
	body, err := rs.Do(ConfigListremotes, nil)
	if err != nil {
		return nil, fmt.Errorf("listremotes: %v", err)
	}
	resp := ListremotesResponse{}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("listremotes parse resp body %s err: %v", removeNewlines(string(body)), err)
	}
	return resp.Remotes, nil
}
