package model

type User struct {
    ID       int64  `db:"id" json:"id"`
    Username string `db:"username" json:"username"`
    Password string `db:"password" json:"-"` // excluded from json responses
    Spv      *int64 `db:"spv" json:"spv,omitempty"`
    Group    string `db:"group" json:"group"`
}
