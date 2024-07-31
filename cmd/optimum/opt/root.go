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
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/fogfish/gurl/v2/http"
	"github.com/fogfish/gurl/x/awsapi"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

// Execute is entry point for cobra cli application
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		e := err.Error()
		fmt.Println(strings.ToUpper(e[:1]) + e[1:])
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&host, "url", "u", "", "host")
	rootCmd.PersistentFlags().StringVarP(&cask, "cask", "c", "", "cask")
}

var (
	host string
	cask string
)

var rootCmd = &cobra.Command{
	Use:   "optimum",
	Short: "tbd",
	Long: `
tbd
	`,
	Run: root,
}

func root(cmd *cobra.Command, args []string) {
	cmd.Help()
}

func stack() (http.Stack, error) {
	aws, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}

	stack := http.New(
		awsapi.WithSignatureV4(aws),
		// http.WithDebugPayload(),
	)

	return stack, nil
}

const IDLE_TIME = 20 * time.Second

func spinner(bar *progressbar.ProgressBar, f func() error) error {
	ch := make(chan bool)

	go func() {
		for {
			select {
			case <-ch:
				return
			default:
				bar.Add(1)
				time.Sleep(40 * time.Millisecond)
			}
		}
	}()

	err := f()

	ch <- false
	bar.Finish()

	return err
}
