package main

import "github.com/satori/go.uuid"

//Link ...
type Link struct {
	Rel    string `json:"rel"`
	Method string `json:"method"`
	Href   string `json:"href"`
}

//FindLinkByRel ...
func FindLinkByRel(rel string, links []Link) LinkResult {
	var result = LinkResult{}
	for _, link := range links {
		if link.Rel == rel {
			result.Found = true
			result.Result = link
			break
		}
	}
	return result
}

//LinkResult ...
type LinkResult struct {
	Found  bool
	Result Link
}

//APIResponse ...
type APIResponse struct {
	Links []Link `json:"_links"`
}

//EndpointDTO ...
type EndpointDTO struct {
	ID          string `json:"id"`
	Destination string `json:"destination"`
	Path        string `json:"path"`
}

//Transform ...
type Transform struct {
	ContentType string
	Transform   string
}

//Endpoint ...
type Endpoint struct {
	ID                uuid.UUID
	Destination       string
	Path              string
	ResponseTransform Transform
}
