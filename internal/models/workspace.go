package models

type Workspace struct {
	Title     string   `json:"title"`
	Desc      string   `json:"desc"`
	Boards    []Board  `json:"boards"`
	Labradors []string `json:"labradors"` //users ID
}
