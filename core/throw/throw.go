// Package throw
// Created by Teocci.
// Author: teocci@yandex.com on 2021-Aug-16
package throw

import (
	"errors"
	"fmt"
)

const (
	errUnableFindSelection = "unable to find %s"
)

func ErrorUnableFindSelection(fqn string) error {
	return errors.New(fmt.Sprintf(errUnableFindSelection, fqn))
}
