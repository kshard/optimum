//
// Copyright (C) 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/kshard/optimum
//

package common

import (
	"context"
	"fmt"
	"time"

	"github.com/fogfish/curie"
	"github.com/kshard/optimum"
)

func AboutList(kind, extension string) string {
	return fmt.Sprintf(`
List all data structure instances. It fetches data structure instances of
type "%s". For each provisioned instance it reports NAME, active VERSION,
UPDATED AT timestamp, instance STATUS, PENDING version if any, and initialization
PARAMS.

  optimum %[1]s list -u $HOST

  NAME      VERSION          UPDATED AT          | STATUS   PENDING          | PARAMS
  example1  NjqOYyOkpMHfg3.6 2024-08-18 10:40:34 | ACTIVE                    | {}
  example2                   2024-08-18 10:38:13 | PENDING  NjqOYyOkpMHfg3.6 | {}

The STATUS reflect both status of the instance and ongoing update operation:
- "UNAVAILABLE" the instance is not ready for use.
- "PENDING" the instance is pending updates, the VERSION is available online. 
- "ACTIVE" the instance is active, all past updates successfully completed.
- "FAILED" PENDING update is failed, the VERSION is available online.
%s
`, kind, extension)
}

// List all data structures of given type
func List(api *optimum.Client, kind string) (err error) {
	seq, err := api.Casks(context.Background(), kind)
	if err != nil {
		return err
	}

	fmt.Printf("%-10s\t%-16s %-19s | %-11s %-16s | %s\n", "NAME", "VERSION", "UPDATED AT", "STATUS", "PENDING", "PARAMS")
	for _, x := range seq.Items {
		fmt.Printf("%-10s\t%-16s %-19s | %-11s %-16s | %s\n", curie.Reference(x.ID), x.Version, x.Updated.Format(time.DateTime), x.Status, x.Pending, x.Opts)
	}

	return nil
}
