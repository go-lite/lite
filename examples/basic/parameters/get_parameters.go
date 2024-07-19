package parameters

type GetReq struct {
	Login  string  `lite:"header=Basic,isauth,scheme=basic,name=Basic"`
	Name   string  `lite:"header=name"`
	Value  *string `lite:"header=value"`
	Params string  `lite:"path=params"`
}
