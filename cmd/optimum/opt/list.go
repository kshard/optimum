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
	"time"

	"github.com/kshard/optimum"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "tbd",
	Long: `
tbd
`,
	SilenceUsage: true,
	RunE:         list,
}

func list(cmd *cobra.Command, args []string) (err error) {
	cli, err := stack()
	if err != nil {
		return err
	}

	api := optimum.New(cli, host)

	seq, err := api.Casks(context.Background(), cask)
	if err != nil {
		return err
	}

	for _, x := range seq.Items {
		fmt.Printf("%10s\t%8s | %s | opts %s\n", x.ID, x.Status, x.Updated.Format(time.DateTime), x.Opts)
	}

	return nil
}
