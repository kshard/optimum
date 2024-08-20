//
// Copyright (C) 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/kshard/optimum
//

package optimum

import (
	"context"

	"github.com/fogfish/curie"
	"github.com/fogfish/gurl/v2/http"
	ƒ "github.com/fogfish/gurl/v2/http/recv"
	ø "github.com/fogfish/gurl/v2/http/send"
	"github.com/fogfish/schemaorg"
)

type Client struct {
	http.Stack

	host ø.Authority
}

func New(stack http.Stack, host string) *Client {
	return &Client{
		Stack: stack,
		host:  ø.Authority(host),
	}
}

func (api *Client) Casks(ctx context.Context, schema string) (*Instances, error) {
	return http.IO[Instances](
		api.WithContext(ctx),
		http.GET(
			ø.URI("%s/ds/%s", api.host, schema),
			ø.Accept.JSON,

			ƒ.Status.OK,
		),
	)
}

func (api *Client) Create(ctx context.Context, cask curie.IRI, opts map[string]any) (*Created, error) {
	return http.IO[Created](
		api.WithContext(ctx),
		http.POST(
			ø.URI("%s/ds/%s", api.host, curie.Prefix(cask)),
			ø.Accept.JSON,
			ø.ContentType.JSON,
			ø.Send(create{
				Name: curie.Reference(cask),
				Opts: opts,
			}),

			ƒ.Status.Accepted,
		),
	)
}

func (api *Client) Commit(ctx context.Context, cask curie.IRI) (*Committed, error) {
	return http.IO[Committed](
		api.WithContext(ctx),
		http.POST(
			ø.URI("%s/ds/%s/%s", api.host, curie.Prefix(cask), curie.Reference(cask)),
			ø.Accept.JSON,
			ø.ContentType.JSON,
			ø.Send(commit{Cursor: "latest"}),

			ƒ.Status.Accepted,
		),
	)

}

func (api *Client) Status(ctx context.Context, job schemaorg.Url) (*JobStatus, error) {
	return http.IO[JobStatus](
		api.WithContext(ctx),
		http.GET(
			ø.URI("%s%s", api.host, ø.Path(job)),
			ø.Accept.JSON,

			ƒ.Status.OK,
		),
	)
}

func (api *Client) Remove(ctx context.Context, cask curie.IRI) error {
	return api.IO(ctx,
		http.DELETE(
			ø.URI("%s/ds/%s/%s", api.host, curie.Prefix(cask), curie.Reference(cask)),
			ø.Accept.JSON,

			ƒ.Status.Accepted,
		),
	)
}

func (api *Client) Query(ctx context.Context, cask curie.IRI, q Query) (*Result, error) {
	return http.IO[Result](
		api.WithContext(ctx),
		http.GET(
			ø.URI("%s/ds/%s/%s", api.host, curie.Prefix(cask), curie.Reference(cask)),
			ø.Accept.JSON,
			ø.ContentType.JSON,
			ø.Send(q),

			ƒ.Status.OK,
		),
	)
}
