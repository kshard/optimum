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

	"github.com/fogfish/curie"
	"github.com/kshard/optimum"
)

func AboutRemove(kind, extension string) string {
	return fmt.Sprintf(`
The command removes "%s" data structure instance. The operation is irreversible and
results in the permanent destruction of all data.
%s
`, kind, extension)
}

func Remove(api *optimum.Client, id curie.IRI) (err error) {
	err = api.Remove(context.Background(), id)
	if err != nil {
		return err
	}

	return nil
}
