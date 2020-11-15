package request
type ReqBookmark struct {
	RepoName string `json:"repo,omitempty" validate:"required"`
}