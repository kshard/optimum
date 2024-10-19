//
// Copyright (C) 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/kshard/optimum
//

package opt

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fogfish/curie"
	"github.com/kshard/optimum"
	"github.com/kshard/optimum/cmd/optimum/encoding"
	"github.com/kshard/optimum/cmd/optimum/opt/common"
	"github.com/kshard/optimum/surface"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

const TYPE_HNSW = "hnsw"

func init() {
	rootCmd.AddCommand(hnswCmd)

	hnswCmd.AddCommand(hnswListCmd)

	hnswCmd.AddCommand(hnswCreateCmd)
	hnswCreateCmd.Flags().StringVarP(&hnswOpts, "json", "j", "", "json config file")

	hnswCmd.AddCommand(hnswCommitCmd)

	hnswCmd.AddCommand(hnswUploadCmd)
	hnswUploadCmd.Flags().IntVar(&hnswUploadBuf, "buf", 4, "upload buffer in MB (default 4MB)")

	hnswCmd.AddCommand(hnswStreamCmd)
	hnswStreamCmd.Flags().IntVar(&hnswChunkSize, "chunk", 100, "streaming chunk size (default 10)")

	hnswCmd.AddCommand(hnswQueryCmd)
	hnswQueryCmd.Flags().StringVarP(&hnswQueryContent, "text", "t", "", "hash to text associated list, useful for debug purposes")

	hnswCmd.AddCommand(hnswRemoveCmd)
}

var (
	hnswOpts         string
	hnswUploadBuf    int
	hnswChunkSize    int
	hnswQueryContent string
)

var hnswCmd = &cobra.Command{
	Use:   "hnsw",
	Short: "Operates `hnsw` data structures.",
	Long: `
The HNSW (Hierarchical Navigable Small World) algorithm is widely applicable in
areas that require efficient nearest neighbor searches, particularly in
high-dimensional spaces. Below are some key areas where HNSW is applicable:

* Text Search and Retrieval: HNSW can be used in search engines and document
retrieval systems to quickly find similar documents on high-dimensional
text embeddings.

* Content-Based Recommendations: HNSW is useful for finding similar items in
recommendation systems, such as finding related products, movies, or music
tracks based on embeddings vectors. It helps quickly locate users or items
with similar behavior patterns.

* Personalized Content: When a system needs to recommend personalized content
(e.g., news articles, blog posts), HNSW can quickly find the most relevant
content based on a user's preferences or behavior.

* Image and Video Retrieval: In tasks like image search or video retrieval,
HNSW is used to find images or frames similar to a given query image,
based on feature vectors extracted from deep learning models.

* Semantic Search: In NLP, HNSW is used to find semantically similar phrases,
sentences, or documents by comparing embeddings generated by text  models.

* Chatbots and Conversational AI: It can be used to match user queries to a set
of predefined responses or intents based on vector similarity.

* Fraud Detection: HNSW can be used to detect anomalies in financial transactions
by identifying transactions that are distant from normal patterns.

* Intrusion Detection: In cybersecurity, it helps to find unusual patterns in
network traffic that might indicate security breaches.
`,
	SilenceUsage: true,
	Run:          hnsw,
}

func hnsw(cmd *cobra.Command, args []string) {
	cmd.Help()
}

//------------------------------------------------------------------------------

var hnswListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all instances of `hnsw` data structure.",
	Long:  common.AboutList(TYPE_HNSW, ""),
	Example: `
optimum hnsw list -u $HOST
optimum hnsw list -u $HOST -r $ROLE
`,
	SilenceUsage: true,
	RunE:         hnswList,
}

func hnswList(cmd *cobra.Command, args []string) (err error) {
	cli, err := stack()
	if err != nil {
		return err
	}

	return common.List(optimum.New(cli, host), TYPE_HNSW)
}

//------------------------------------------------------------------------------

var hnswCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create new instance of `hnsw` data structure.",
	Long: common.AboutCreate("hnsw", `
The algorithm "hnsw" is an efficient and scalable method for approximate nearest
neighbor search in high-dimensional spaces.

Config algorithm through primary parameters: 
  - "M" and "M0" controls the maximum number of connections per node, balancing
    between memory usage and search efficiency.  M0 defines the connection
		density on the graph's base layer, while M regulates it on the intermediate
		layers.

  - "efConstruction" determines the number of candidate nodes evaluated during
    graph construction, influencing both the construction time and the accuracy
    of the graph.

  - "surface" is vector distance function.

Example configuration and default values:	
  {
    "m":  8,                // number in range of [4, 1024]
    "m0": 64,               // number in range of [4, 1024]
    "efConstruction": 200,  // number in range of [200, 1000]
    "surface": "cosine"     // enum {"cosine", "euclidean"}
  }

`),
	Example: `
optimum hnsw create -u $HOST -n example -j path/to/config.json
optimum hnsw create -u $HOST -r $ROLE -n example -j path/to/config.json
`,
	SilenceUsage: true,
	RunE:         hnswCreate,
}

func hnswCreate(cmd *cobra.Command, args []string) (err error) {
	cli, err := stack()
	if err != nil {
		return err
	}

	return common.Create(optimum.New(cli, host), curie.New("%s:%s", TYPE_HNSW, name), hnswOpts)
}

//------------------------------------------------------------------------------

var hnswCommitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Commit earlier uploaded datasets into `hnsw` instance.",
	Long:  common.AboutCommit(TYPE_HNSW, ""),
	Example: `
optimum hnsw commit -u $HOST -n example
optimum hnsw commit -u $HOST -r $ROLE -n example
`,
	SilenceUsage: true,
	RunE:         hnswCommit,
}

func hnswCommit(cmd *cobra.Command, args []string) (err error) {
	cli, err := stack()
	if err != nil {
		return err
	}

	return common.Commit(optimum.New(cli, host), curie.New("%s:%s", TYPE_HNSW, name))
}

//------------------------------------------------------------------------------

var hnswUploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload `hnsw` datasets.",
	Long: `
Upload "hnsw" dataset to server. It accepts only textual format to represent
embedding vectors. Each line of the file should start with unique key, followed
by the corresponding vector. The unique key length should not exceeding 32 bytes:

  example_key_a 0.24116 ... -0.26098 -0.0079604
	example_key_b 0.34601 ... -0.66865 -0.0486001

We recommend using sha1, uuid or https://github.com/fogfish/guid as unique key.
The format allows hexadecimal encoding for keys, if it starts with "0x" prefix.    

  0xd857f9dc157c28e8e07c569c5992dee4f3486b4c -0.097231 ... -0.001681 0.154977
  0xaeb3e05ab60520cd947455f2130d6cf1f6103243 -0.008007 ... -0.098503 0.057056

`,
	Example: `
optimum hnsw upload -u $HOST -n example path/to/data.txt
optimum hnsw upload -u $HOST -r $ROLE -n example path/to/data.txt
`,
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	RunE:         hnswUpload,
}

func hnswUpload(cmd *cobra.Command, args []string) (err error) {
	fd, err := os.Open(args[0])
	if err != nil {
		return err
	}
	defer fd.Close()

	fi, err := fd.Stat()
	if err != nil {
		return err
	}

	cli, err := stack()
	if err != nil {
		return err
	}

	stream := surface.NewWriter(cli, host, curie.New("%s:%s", TYPE_HNSW, name), hnswUploadBuf*1024*1024)

	r := io.TeeReader(fd,
		progressbar.DefaultBytes(
			fi.Size(),
			"==> uploading",
		),
	)

	scanner := encoding.New(r)
	for scanner.Scan() {
		err := stream.Write(context.Background(),
			surface.Vector{
				UniqueKey: scanner.UniqueKey(),
				Vector:    scanner.Vector(),
			},
		)
		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if err := stream.Sync(context.Background()); err != nil {
		return err
	}

	return nil
}

//------------------------------------------------------------------------------

var hnswStreamCmd = &cobra.Command{
	Use:   "stream",
	Short: "Stream `hnsw` datasets.",
	Long: `
Stream "hnsw" dataset to server. It accepts only textual format to represent
embedding vectors. Each line of the file should start with unique key, followed
by the corresponding vector. The unique key length should not exceeding 32 bytes:

  example_key_a 0.24116 ... -0.26098 -0.0079604
	example_key_b 0.34601 ... -0.66865 -0.0486001

We recommend using sha1, uuid or https://github.com/fogfish/guid as unique key.
The format allows hexadecimal encoding for keys, if it starts with "0x" prefix.    

  0xd857f9dc157c28e8e07c569c5992dee4f3486b4c -0.097231 ... -0.001681 0.154977
  0xaeb3e05ab60520cd947455f2130d6cf1f6103243 -0.008007 ... -0.098503 0.057056

`,
	Example: `
optimum hnsw stream -u $HOST -n example path/to/data.txt
optimum hnsw stream -u $HOST -r $ROLE -n example path/to/data.txt
`,
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	RunE:         hnswStream,
}

func hnswStream(cmd *cobra.Command, args []string) (err error) {
	fd, err := os.Open(args[0])
	if err != nil {
		return err
	}
	defer fd.Close()

	fi, err := fd.Stat()
	if err != nil {
		return err
	}

	cli, err := stack()
	if err != nil {
		return err
	}

	// sentences
	// surface
	api := surface.New(cli, host)

	r := io.TeeReader(fd,
		progressbar.DefaultBytes(
			fi.Size(),
			"==> uploading",
		),
	)

	scanner := encoding.New(r)
	for scanner.Scan() {
		bag := make([]surface.Vector, 0)
		for i, has := 0, true; i < hnswChunkSize && has; i, has = i+1, scanner.Scan() {
			bag = append(bag, surface.Vector{
				UniqueKey: scanner.UniqueKey(),
				Vector:    scanner.Vector(),
			})
		}

		if len(bag) > 0 {
			err := api.Write(context.Background(), curie.New("%s:%s", TYPE_HNSW, name), bag)
			if err != nil {
				return err
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

//------------------------------------------------------------------------------

var hnswQueryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query instance of `hnsw` data structure.",
	Long: `
Query "hnsw" data structure instance. It accepts textual format as input, where
each line is embedding vector to query. Each line of the file should start with
identity of query, followed by the corresponding vector: 

  example_query_a 0.24116 ... -0.26098 -0.0079604
	example_query_b 0.34601 ... -0.66865 -0.0486001

The file format is identical to the upload and can be re-used as is.
`,
	Example: `
optimum hnsw query -u $HOST -n example path/to/query.txt
optimum hnsw query -u $HOST -r $ROLE -n example path/to/query.txt
optimum hnsw query -u $HOST -r $ROLE -n example -t path/to/text-map.txt path/to/query.txt
`,
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	RunE:         hnswQuery,
}

func hnswQuery(cmd *cobra.Command, args []string) (err error) {
	hashmap := hnswTextHashMap()

	fd, err := os.Open(args[0])
	if err != nil {
		return err
	}
	defer fd.Close()

	cli, err := stack()
	if err != nil {
		return err
	}

	api := surface.New(cli, host)

	scanner := encoding.New(fd)
	for scanner.Scan() {
		query := surface.Query{Query: scanner.Vector()}
		rs, err := api.Query(context.Background(), curie.New("%s:%s", TYPE_HNSW, name), query)
		if err != nil {
			return err
		}

		id := fmt.Sprintf("0x%x", scanner.UniqueKey())
		fmt.Printf("Query %s (took %s) | %s (vsn %s, size %d)\n", id, rs.Took, rs.Source.Cask, rs.Source.Version, rs.Source.Size)
		for _, hit := range rs.Hits {
			hid := fmt.Sprintf("0x%x", hit.UniqueKey)
			fmt.Printf("  %f : %32s \n", hit.Rank, hid)
		}

		if hashmap != nil {
			fmt.Printf("\n\nQuery (took %s) > %s\n", rs.Took, hnswTextValue(hashmap, id))
			for _, hit := range rs.Hits {
				hid := fmt.Sprintf("0x%x", hit.UniqueKey)
				fmt.Printf("  %f : %s\n", hit.Rank, hnswTextValue(hashmap, hid))
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func hnswTextHashMap() map[string]string {
	if hnswQueryContent == "" {
		return nil
	}

	fd, err := os.Open(hnswQueryContent)
	if err != nil {
		return nil
	}
	defer fd.Close()

	hashmap := map[string]string{}
	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		seq := strings.SplitN(scanner.Text(), " ", 2)
		hashmap[seq[0]] = seq[1]
	}

	return hashmap
}

func hnswTextValue(hashmap map[string]string, key string) string {
	if val, has := hashmap[key]; has {
		return val
	}

	return key
}

//------------------------------------------------------------------------------

var hnswRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove instance of `hnsw` data structure.",
	Long:  common.AboutRemove(TYPE_HNSW, ""),
	Example: `
optimum hnsw commit -u $HOST -n example
optimum hnsw commit -u $HOST -r $ROLE -n example
`,
	SilenceUsage: true,
	RunE:         hnswRemove,
}

func hnswRemove(cmd *cobra.Command, args []string) (err error) {
	cli, err := stack()
	if err != nil {
		return err
	}

	return common.Remove(optimum.New(cli, host), curie.New("%s:%s", TYPE_HNSW, name))
}
