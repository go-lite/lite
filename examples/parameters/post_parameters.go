package parameters

import "github.com/disco07/lite-fiber/codec"

type CreateBody struct {
	FirstName string  `json:"first_name"`
	LastName  *string `json:"last_name"`
}

type Params struct {
	ID uint64 `params:"id"`
}

type ReqHeader struct {
	Authorization *string `reqHeader:"Authorization" header:"Authorization"`
}

type CreateReq struct {
	Header codec.AsHeader[ReqHeader]
	Params codec.AsPathParam[Params]
	Body   codec.AsJSON[CreateBody]
}
