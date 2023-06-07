// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package generated

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID        pgtype.UUID        `json:"id"`
	Username  string             `json:"username"`
	Password  string             `json:"password"`
	Avatar    pgtype.Text        `json:"avatar"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}