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
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/fogfish/curie"
	"github.com/kshard/optimum"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&fopts, "json", "j", "", "json config file")
}

var (
	fopts string
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create new data structure.",
	Long: `
Create new cask (data structure) on the cloud. Specify data structure
algorithm (class) and a unique name along with configuration properties.
The configuration properties are supplied via json file.

See below the list of supported algorithms:


1. Hierarchical Navigable Small World

The algorithm "hnsw" is an efficient and scalable method for approximate nearest
neighbor search in high-dimensional spaces.

Config algorithm through primary parameters: 
  - "M" and "M0" controls the maximum number of connections per node, balancing
    between memory usage and search efficiency.

  - "efConstruction" determines the number of candidate nodes evaluated during
    graph construction, influencing both the construction time and the accuracy
    of the graph.

  - "surface" is vector distance function.

Example configuration:	
  {
    "m":  8,                // number in range of [4, 1024]
    "m0": 64,               // number in range of [4, 1024]
    "efConstruction": 200,  // number in range of [200, 1000]
    "surface": "cosine"     // enum {"cosine", "euclidean"}
  }
`,
	Example: `
optimum create -u $HOST -c class:cask -j path/to/config.json
optimum create -u $HOST -r $ROLE -c class:cask -j path/to/config.json
`,
	SilenceUsage: true,
	RunE:         create,
}

func create(cmd *cobra.Command, args []string) (err error) {
	opts := map[string]any{}

	if fopts != "" {
		b, err := os.ReadFile(fopts)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(b, &opts); err != nil {
			return err
		}
	}

	cli, err := stack()
	if err != nil {
		return err
	}

	api := optimum.New(cli, host)

	receipt, err := api.Create(context.Background(), curie.IRI(cask), opts)
	if err != nil {
		return err
	}

	bar := progressbar.NewOptions(-1,
		progressbar.OptionSpinnerType(14),
		progressbar.OptionSetDescription(
			fmt.Sprintf("%s | %s ... opts: %+v", cask, "CREATING", opts),
		),
	)

	return spinner(bar, func() error {
		for {
			time.Sleep(IDLE_TIME)

			status, err := api.Status(context.Background(), receipt.Job)
			if err != nil {
				return err
			}

			bar.Describe(fmt.Sprintf("%s | %s ...", cask, status.Status))
			if status.Status == "SUCCEEDED" || status.Status == "FAILED" {
				return nil
			}
		}
	})
}
