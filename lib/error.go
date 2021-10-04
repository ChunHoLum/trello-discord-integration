package lib

import (
	"context"
)

func IsCanceled(err error) bool {
	return err == context.Canceled
}

func IsDeadline(err error) bool {
	return err == context.DeadlineExceeded
}
