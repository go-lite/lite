package parameters

type GetReq struct {
	Login  string `lite:"header=Basic,isauth,scheme=basic,name=Basic"`
	Params string `lite:"path=name"`
}
