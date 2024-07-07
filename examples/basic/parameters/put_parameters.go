package parameters

type PutReq struct {
	ID   uint64  `lite:"path=id"`
	Body PutBody `lite:"req=body"`
}

type PutBody struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
