package returns

type CreateResponse struct {
	ID        uint64  `lite:"resp:json" json:"id"`
	FirstName string  `lite:"resp:json" json:"fist_name"`
	LastName  *string `lite:"resp:json" json:"last_name"`
}
