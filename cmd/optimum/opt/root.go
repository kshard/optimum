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

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
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
	rootCmd.PersistentFlags().StringVarP(&host, "url", "u", "", "url to remote data structure management server")
	rootCmd.PersistentFlags().StringVarP(&cask, "cask", "c", "", "cask identity in the format (class:name), class defines data structure algorithm, followed by a unique name after :.")
	rootCmd.PersistentFlags().StringVarP(&role, "role", "r", "", "access identity, ARN of AWS IAM Role")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug output")
}

var (
	host  string
	cask  string
	role  string
	debug bool
)

var rootCmd = &cobra.Command{
	Use:   "optimum",
	Short: "client for managing cloud data structures",
	Long: `
The command line client for managing cloud data structures. The data structure
is a collection, referred to as casks. Each cask is implemented based on
a specific data structure algorithm (class) and is assigned a unique name
along with configuration properties. This utility helps cask management on your
behalf.

The command line utility requires access to remote server that provisions and
operates data structures for you. Contact your provided for details.

It is recommended to config environment variables for client usage:

  export HOST=https://example.com
  export ROLE=arn:aws:iam::000000000000:role/example-access-role

	`,
	Run: root,
}

func root(cmd *cobra.Command, args []string) {
	cmd.Help()
}

func stack() (http.Stack, error) {
	opts := []http.Config{}

	if debug {
		opts = append(opts, http.WithDebugPayload())
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}

	if role == "" {
		opts = append(opts, awsapi.WithSignatureV4(cfg))
	} else {
		assumed, err := config.LoadDefaultConfig(context.Background(),
			config.WithCredentialsProvider(
				aws.NewCredentialsCache(
					stscreds.NewAssumeRoleProvider(sts.NewFromConfig(cfg), role),
				),
			),
		)
		if err != nil {
			return nil, err
		}

		opts = append(opts, awsapi.WithSignatureV4(assumed))
	}

	return http.New(opts...), nil
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
