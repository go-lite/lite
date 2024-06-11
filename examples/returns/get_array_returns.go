package returns

type GetArrayReturnsResponse = []Ret

type Ret struct {
	Message string `json:"message"`
	Embed   Embed  `json:"embed"`
}

type Embed struct {
	Key        string     `json:"key"`
	ValueEmbed ValueEmbed `json:"value"`
	Others     []*string  `json:"others"`
	OtherEmbed OtherEmbed `json:"other_embed"`
}

type ValueEmbed = *string

type OtherEmbed struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
