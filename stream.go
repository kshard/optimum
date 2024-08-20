//
// Copyright (C) 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/kshard/optimum
//

package optimum

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha1"

	"github.com/fogfish/curie"
	"github.com/fogfish/gurl/v2/http"
	ƒ "github.com/fogfish/gurl/v2/http/recv"
	ø "github.com/fogfish/gurl/v2/http/send"
	"github.com/kshard/wreck"
)

// Stream client
type Stream struct {
	http.Stack

	host ø.Authority
	cask curie.IRI

	chunk int
	buf   *bytes.Buffer
	zip   *gzip.Writer
	seq   *wreck.Writer[float32]
}

func NewStream(stack http.Stack, host string, cask curie.IRI, chunk int) *Stream {
	stream := &Stream{
		Stack: stack,
		host:  ø.Authority(host),
		cask:  cask,
		chunk: chunk,
	}

	stream.buf = &bytes.Buffer{}
	stream.zip = gzip.NewWriter(stream.buf)
	stream.seq = wreck.NewWriter[float32](stream.zip)

	return stream
}

// Write vector
func (stream *Stream) Write(ctx context.Context, v Vector) error {
	if err := stream.seq.Write(v.UniqueKey, v.SortKey, v.Vec); err != nil {
		return err
	}

	if stream.buf.Len() >= stream.chunk {
		return stream.Sync(ctx)
	}

	return nil
}

// Sync local cache
func (stream *Stream) Sync(ctx context.Context) (err error) {
	if err := stream.zip.Close(); err != nil {
		return err
	}

	if stream.buf.Len() == 0 {
		return nil
	}

	return stream.Stack.IO(ctx,
		http.PUT(
			ø.URI("%s/ds/%s/%s", stream.host, curie.Prefix(stream.cask), curie.Reference(stream.cask)),
			ø.Accept.JSON,
			ø.ContentType.Set("application/octet-stream"),
			ø.Send(stream.buf),

			ƒ.Status.Accepted,
			func(ctx *http.Context) error {
				stream.buf.Reset()
				stream.zip.Reset(stream.buf)
				return nil
			},
		),
	)
}

// Textual stream client
type TextStream struct {
	api    Embeddings
	stream *Stream
}

type Embeddings interface {
	Embedding(ctx context.Context, text string) ([]float32, error)
}

func NewTextStream(api Embeddings, stream *Stream) *TextStream {
	return &TextStream{
		api:    api,
		stream: stream,
	}
}

func (stream *TextStream) Write(ctx context.Context, text string) error {
	vec, err := stream.api.Embedding(ctx, text)
	if err != nil {
		return err
	}

	hash := sha1.New()
	hash.Write([]byte(text))
	uniqueKey := hash.Sum(nil)

	v := Vector{
		UniqueKey: uniqueKey,
		Vec:       vec,
	}

	if err := stream.stream.Write(ctx, v); err != nil {
		return err
	}

	return nil
}

func (stream *TextStream) Sync(ctx context.Context) (err error) {
	return stream.stream.Sync(ctx)
}
