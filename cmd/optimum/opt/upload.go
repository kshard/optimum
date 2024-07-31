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
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/fogfish/curie"
	"github.com/kshard/optimum"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(uploadCmd)
	uploadCmd.Flags().IntVar(&chunk, "chunk", 1000, "chunk size")
}

var (
	chunk int
)

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "tbd",
	Long: `
tbd
`,
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	RunE:         upload,
}

func upload(cmd *cobra.Command, args []string) (err error) {
	cli, err := stack()
	if err != nil {
		return err
	}

	fd, err := os.Open(args[0])
	if err != nil {
		return err
	}
	defer fd.Close()

	fi, err := fd.Stat()
	if err != nil {
		return err
	}

	bar := progressbar.DefaultBytes(
		fi.Size(),
		"==> uploading",
	)

	api := optimum.NewStream(cli, host, curie.IRI(cask), chunk)

	return scan(io.TeeReader(fd, bar),
		func(key string, vec []float32) error {
			if len(key) > 31 {
				key = key[:31]
			}

			return api.Write(context.Background(),
				optimum.Vector{
					Vec:       vec,
					UniqueKey: []byte(key),
				},
			)
		},
	)
}

func scan(r io.Reader, f func(string, []float32) error) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		txt := scanner.Text()
		seq := strings.Split(txt, " ")

		vec := make([]float32, len(seq)-1)
		for i := 1; i < len(seq); i++ {
			v, err := strconv.ParseFloat(seq[i], 32)
			if err != nil {
				return err
			}
			vec[i-1] = float32(v)
		}

		key := seq[0]
		if err := f(key, vec); err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
