package model

type Approval struct {
    ID          int64  `db:"id" json:"id"`
    Transaction string `db:"transaction" json:"transaction"`
    Requester   int64  `db:"requester" json:"requester"`
    Approver    int64  `db:"approver" json:"approver"`
}
