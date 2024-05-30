package swagger

type Route struct {
	Method           string
	Path             string
	Name             string
	Tags             []string
	RequestsResponse RequestResponse
}

type RequestResponse struct {
	Request  Request
	Response Response
}

type Security struct {
	Type   string
	Name   string
	Scheme string
}

type Request struct {
	Headers         []Header
	PathParameters  []PathParameter
	QueryParameters []QueryParameter
	Security        []Security
	Body            RequestBody
}

type Response struct {
	Headers []Header
	Body    ResponseBody
}

type Header struct {
	Key   string
	Value Value
}

type PathParameter struct {
	Key   string
	Value Value
}

type QueryParameter struct {
	Key   string
	Value Value
}

type RequestBody map[string]Value

type ResponseBody map[int]map[string]Value

type ValueType string

const (
	ValueTypeString ValueType = "string"
	ValueTypeInt    ValueType = "int"
	ValueTypeUint   ValueType = "uint"
	ValueTypeFloat  ValueType = "float"
	ValueTypeBool   ValueType = "bool"
	ValueTypeObject ValueType = "object"
	ValueTypeArray  ValueType = "array"
	ValueTypeMap    ValueType = "map"
	ValueTypeAny    ValueType = "any"
)

type ValueFormat string

const (
	ValueFormatDateTime ValueFormat = "date-time"
	ValueFormatBinary   ValueFormat = "binary"
)

type Value struct {
	Name        string
	Type        ValueType
	Nullable    bool
	Format      ValueFormat
	Properties  map[string]Value
	NoExplode   bool
	Items       *Value
	Keys        *Value
	Values      *Value
	Enum        []string
	XMLName     string
	ContentType []string
}
