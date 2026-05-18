package usecase

import "context"

type AlwaysExistsChecker struct{}

func (AlwaysExistsChecker) Exists(context.Context, string) (bool, error) {
	return true, nil
}