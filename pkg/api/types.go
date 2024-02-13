package api

type Links map[string]Link

type Embedded map[string]interface{}

type LinkOpts struct {
	Links Links `json:"_links"`
}

type Collection struct {
	LinkOpts
	Metadata CollectionOpts `json:"_metadata"`
	Embedded Embedded       `json:"_embedded"`
}

type CollectionOpts struct {
	Offset int   `json:"offset"`
	Limit  int   `json:"limit"`
	Total  int64 `json:"total"`
}

type Link struct {
	Href string `json:"href"`
}
