package returns

type GetArrayReturnsResponse = []Ret

type Ret struct {
	Message  string                 `json:"message"`
	Embed    Embed                  `json:"embed"`
	Map      map[string]string      `json:"map"`
	OtherMap map[string]OtherEmbed2 `json:"other_map"`
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

type OtherEmbed2 struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
