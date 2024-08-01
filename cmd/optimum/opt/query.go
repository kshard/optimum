//
// Copyright (C) 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/kshard/optimum
//

package opt

import (
	"context"
	"fmt"
	"os"

	"github.com/fogfish/curie"
	"github.com/kshard/optimum"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(queryCmd)
}

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query data structure",
	Long: `
Query data structure. Use command for basic operations only. The client supports
the following query formats, which vary depending on the data structure
algorithms used: 


1. Hierarchical Navigable Small World

The query is equivalent to data format used for writing. Textual formats are
required to represent embedding vectors, with each line consisting of a key
(query ID) and a vector.

  query -0.37604 0.24116 ... -0.26098 -0.0079604

`,
	Example: `
optimum query -u $HOST -c class:cask path/to/query.txt
`,
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	RunE:         query,
}

func query(cmd *cobra.Command, args []string) (err error) {
	cli, err := stack()
	if err != nil {
		return err
	}

	fd, err := os.Open(args[0])
	if err != nil {
		return err
	}
	defer fd.Close()

	api := optimum.New(cli, host)

	return scan(fd,
		func(key string, vec []float32) error {
			r, err := api.Query(context.Background(), curie.IRI(cask),
				optimum.Query{Query: vec},
			)
			if err != nil {
				return err
			}

			fmt.Printf("Query %s (took %s) | %s (vsn %s, size %d)\n", key, r.Took, r.Version.Cask, r.Version.Version, r.Version.Size)
			for _, hit := range r.Hits {
				fmt.Printf("  %32s : %f | % 0x\n", string(hit.UniqueKey), hit.Rank, hit.UniqueKey)
			}

			return nil
		},
	)
}
