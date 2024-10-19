//
// Copyright (C) 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/kshard/optimum
//

package sentences

import (
	"time"

	"github.com/fogfish/schemaorg"
	"github.com/kshard/optimum"
)

// Sentence defines textual content
type Sentence struct {
	// Short text block.
	Text schemaorg.Text `json:"text,omitempty"`

	// URL of the original CreativeWork from which this text block is derived.
	IsPartOf schemaorg.IsPartOf `json:"isPartOf,omitempty"`

	// Headline(s) of the text block.
	Headline []schemaorg.Headline `json:"headline,omitempty"`

	// Relevant Keywords for the text block.
	Keywords []schemaorg.Keywords `json:"keywords,omitempty"`

	// External URIs associated with the text block.
	//
	// The purpose lacks a defined product type for representing the URI and its
	// human-readable format. The client selects own format.
	//
	// We recommend:
	//	- URIs with Fragments is recommended to represent url and human readable name
	//	<http://example.com/article123#An_Example_Article>
	//
	//	- Markdown format
	//	[An Example Article](http://example.com/article123)
	//
	//	- HTML Anchor Tag
	//	<a href="http://example.com/article123">An Example Article</a>
	//
	//	- URI with Label Convention (RDF-style text)
	//	<http://example.com/article123> "An Example Article"
	Links []schemaorg.RelatedLink `json:"links,omitempty"`
}

// Query textual corpus
type Query struct {
	K        int     `json:"k,omitempty"`
	EfSearch int     `json:"efSearch,omitempty"`
	Distance float32 `json:"distance,omitempty"`
	Text     string  `json:"text"`
}

// Results from the query
type Result struct {
	Took   time.Duration  `json:"took,omitempty"`
	Source optimum.Source `json:"source,omitempty"`
	Hits   []Hit          `json:"hits,omitempty"`
}

type Hit struct {
	Sentence
	Rank float32 `json:"rank,omitempty"`
}
