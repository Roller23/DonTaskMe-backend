package model

type WorkspaceRequest struct {
	Token     string   `json:"token"`
	Title     string   `json:"title"`
	Desc      string   `json:"desc"`
	Boards    []Board  `json:"boards,omitempty"`
	Labradors []string `json:"labradors,omitempty"`
}

type Workspace struct {
	Title     string   `json:"title"`
	Desc      string   `json:"desc"`
	Boards    []Board  `json:"boards"`
	Labradors []string `json:"labradors"` //users ID
}
