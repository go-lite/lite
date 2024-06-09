package parameters

type CreateBody struct {
	FirstName string  `json:"first_name"`
	LastName  *string `json:"last_name"`
}

type Params struct {
	ID uint64 `lite:"path=id"`
}

type ReqHeader struct {
	Authorization *string `lite:"header=Authorization,isauth,scheme=bearer"`
}

type CreateReq struct {
	Header ReqHeader
	Params Params
	Body   CreateBody `lite:"req=body"`
}
