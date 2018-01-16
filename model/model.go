package model

/*
	Refs:
	- https://github.com/harvest-platform/concept-registry
	- https://github.com/harvest-platform/esevaluator
	- https://github.com/SpaceHexagon/media-scraper
	- https://github.com/Akagi201/jsonbeat
*/

type UserVisit struct {
	Date               int64
	User_id            int
	Session_id         string
	Page_id            int
	Action_time        int64
	Search_keyword     string
	Click_category_id  int
	Click_product_id   int
	Order_category_ids string
	Order_product_ids  string
	Pay_category_ids   string
	Pay_product_ids    string
}

// Tree represents a GitHub tree.
type Tree struct {
	SHA     *string     `json:"sha,omitempty"`
	Entries []TreeEntry `json:"tree,omitempty"`
}

func (t Tree) String() string {
	return Stringify(t)
}

// TreeEntry represents the contents of a tree structure. TreeEntry can represent either a blob, a commit (in the case of a submodule), or another tree.
type TreeEntry struct {
	SHA     *string `json:"sha,omitempty"`
	Path    *string `json:"path,omitempty"`
	Mode    *string `json:"mode,omitempty"`
	Type    *string `json:"type,omitempty"`
	Size    *int    `json:"size,omitempty"`
	Content *string `json:"content,omitempty"`
	URL     *string `json:"url,omitempty"`
}

func (t TreeEntry) String() string {
	return Stringify(t)
}
