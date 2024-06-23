package parameters

type CreateBody struct {
	FirstName string  `json:"first_name"`
	LastName  *string `json:"last_name"`
}

type CreateReq struct {
	Authorization *string    `lite:"header=Authorization,isauth,scheme=bearer"`
	ID            uint64     `lite:"path=id"`
	Body          CreateBody `lite:"req=body"`
}
