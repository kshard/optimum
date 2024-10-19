//
// Copyright (C) 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/kshard/optimum
//

package sentences

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"

	"github.com/fogfish/curie"
	"github.com/fogfish/gurl/v2/http"
	ƒ "github.com/fogfish/gurl/v2/http/recv"
	ø "github.com/fogfish/gurl/v2/http/send"
)

// Client for streaming to Graph-based Nearest Neighbor Search Algorithms
type Writer struct {
	http.Stack

	host ø.Authority
	cask curie.IRI

	chunk int
	buf   bytes.Buffer
	zip   *gzip.Writer
	seq   *json.Encoder
}

// Creates the client for streaming to Graph-based Nearest Neighbor Search Algorithms
func NewWriter(stack http.Stack, host string, cask curie.IRI, chunk int) *Writer {
	stream := &Writer{
		Stack: stack,
		host:  ø.Authority(host),
		cask:  cask,
		chunk: chunk,
	}
	stream.reset()

	return stream
}

func (stream *Writer) reset() {
	stream.buf.Reset()
	stream.zip = gzip.NewWriter(&stream.buf)
	stream.seq = json.NewEncoder(stream.zip)
}

// Write vector
func (stream *Writer) Write(ctx context.Context, v Sentence) error {
	if err := stream.seq.Encode(v); err != nil {
		return err
	}

	if stream.buf.Len() >= stream.chunk {
		return stream.Sync(ctx)
	}

	return nil
}

// Sync local cache
func (stream *Writer) Sync(ctx context.Context) (err error) {
	defer stream.reset()

	if err := stream.zip.Close(); err != nil {
		return err
	}

	if stream.buf.Len() == 0 {
		return nil
	}

	return stream.Stack.IO(ctx,
		http.POST(
			ø.URI("%s/ds/%s/%s/objects", stream.host, curie.Prefix(stream.cask), curie.Reference(stream.cask)),
			ø.Accept.JSON,
			ø.ContentType.JSON,
			ø.Send(struct {
				V []byte `json:"object"`
			}{
				V: stream.buf.Bytes(),
			}),

			ƒ.Status.Accepted,
		),
	)
}
