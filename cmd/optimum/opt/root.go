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
	"os/user"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/fogfish/gurl/v2/http"
	"github.com/fogfish/gurl/x/awsapi"
	"github.com/jdxcode/netrc"
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
	rootCmd.PersistentFlags().StringVarP(&name, "name", "n", "", "unique name of data structure, use only alpha-numeric symbols.")
	rootCmd.PersistentFlags().StringVarP(&role, "role", "r", "", "access identity, ARN of AWS IAM Role")
	rootCmd.PersistentFlags().StringVarP(&exid, "external-id", "e", "", "ExternalID associated with the role")
	rootCmd.PersistentFlags().StringVarP(&profile, "profile", "p", "", "the access profile at ~/.aws/config")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug output")
}

var (
	host    string
	name    string
	role    string
	exid    string
	profile string
	debug   bool
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

//------------------------------------------------------------------------------

func stack() (http.Stack, error) {
	if role != "" {
		return stackFromRole()
	}

	if profile != "" {
		return stackFromProfile()
	}

	return stackFromNetRC()
}

func stackFromConfig(cfg aws.Config) (http.Stack, error) {
	opts := []http.Option{awsapi.WithSignatureV4(cfg)}
	if debug {
		opts = append(opts, http.WithDebugPayload)
	}

	return http.New(opts...), nil
}

func stackDefault() (http.Stack, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}

	return stackFromConfig(cfg)
}

func stackFromProfile() (http.Stack, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithSharedConfigProfile(profile)) // config.WithClientLogMode(aws.LogRequestWithBody|aws.LogResponseWithBody),

	if err != nil {
		return nil, err
	}

	return stackFromConfig(cfg)
}

func stackFromRole() (http.Stack, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}

	assumed, err := config.LoadDefaultConfig(context.Background(),
		config.WithCredentialsProvider(
			aws.NewCredentialsCache(
				stscreds.NewAssumeRoleProvider(sts.NewFromConfig(cfg), role,
					func(aro *stscreds.AssumeRoleOptions) {
						if exid != "" {
							aro.ExternalID = aws.String(exid)
						}
					},
				),
			),
		),
	)
	if err != nil {
		return nil, err
	}

	return stackFromConfig(assumed)
}

func stackFromNetRC() (http.Stack, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	file := filepath.Join(usr.HomeDir, ".netrc")
	n, err := netrc.Parse(file)
	if err != nil {
		return nil, err
	}

	machine := n.Machine("optimum")
	if machine == nil {
		return stackDefault()
	}

	// .netrc defines default url & datastore name
	if len(host) == 0 {
		host = machine.Get("host")
	}

	// Using AWS profile
	profile = machine.Get("profile")
	if len(profile) == 0 {
		return stackDefault()
	}

	return stackFromProfile()
}
