//
// Copyright (C) 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/kshard/optimum
//

package sentences

import (
	"context"

	"github.com/fogfish/curie"
	"github.com/fogfish/gurl/v2/http"
	ƒ "github.com/fogfish/gurl/v2/http/recv"
	ø "github.com/fogfish/gurl/v2/http/send"
)

// Client for reading/writing natural language text and searching for nearest neighbor.
type Client struct {
	http.Stack

	host ø.Authority
}

// Creates the client for reading/writing natural language text and searching for nearest neighbor.
func New(stack http.Stack, host string) *Client {
	return &Client{
		Stack: stack,
		host:  ø.Authority(host),
	}
}

// Write the sentence
func (api *Client) Write(ctx context.Context, cask curie.IRI, bag []Sentence) error {
	if len(bag) == 0 {
		return nil
	}

	return api.IO(ctx,
		http.POST(
			ø.URI("%s/ds/%s/%s/object", api.host, curie.Prefix(cask), curie.Reference(cask)),
			ø.Accept.JSON,
			ø.ContentType.JSON,
			ø.Send(struct {
				V []Sentence `json:"object"`
			}{
				V: bag,
			}),

			ƒ.Status.Accepted,
		),
	)
}

// Query nearest neighbor text to the given sample.
func (api *Client) Query(ctx context.Context, cask curie.IRI, q Query) (*Result, error) {
	return http.IO[Result](
		api.WithContext(ctx),
		http.GET(
			ø.URI("%s/ds/%s/%s", api.host, curie.Prefix(cask), curie.Reference(cask)),
			ø.Accept.JSON,
			ø.ContentType.JSON,
			ø.Send(struct {
				Q Query `json:"query"`
			}{
				Q: q,
			}),

			ƒ.Status.OK,
		),
	)
}
