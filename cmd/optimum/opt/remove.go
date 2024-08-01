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

	"github.com/fogfish/curie"
	"github.com/kshard/optimum"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(removeCmd)
}

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove data collection.",
	Long: `
Permanently remove data collection. The operation is irreversible and results
in the permanent destruction of all data with in the collection.
`,
	Example: `
optimum remove -u $HOST -c class:cask
`,
	SilenceUsage: true,
	RunE:         remove,
}

func remove(cmd *cobra.Command, args []string) (err error) {
	cli, err := stack()
	if err != nil {
		return err
	}

	api := optimum.New(cli, host)

	err = api.Remove(context.Background(), curie.IRI(cask))
	if err != nil {
		return err
	}

	return nil
}
