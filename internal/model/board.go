package model

type Board struct {
	Title       string   `json:"title"`
	WorkspaceID string   `json:"workspace_id"`
	OwnerID     string   `json:"owner_id"`
	Labradors   []string `json:"labradors"`
}
