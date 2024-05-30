package parameters

import "github.com/disco07/lite-fiber/codec"

type GetParams struct {
	Name string `params:"name"`
}

type Basic struct {
	Basic string `reqHeader:"Basic" header:"Authorization" scheme:"basic" name:"Basic"`
}

type GetReq struct {
	Login  codec.AsHeader[Basic]
	Params codec.AsPathParam[GetParams]
}
