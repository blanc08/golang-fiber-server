// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0

package database

import (
	"context"
)

type Querier interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	GetUser(ctx context.Context, email string) (User, error)
}

var _ Querier = (*Queries)(nil)