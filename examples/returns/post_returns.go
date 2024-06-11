package returns

type CreateResponse struct {
	ID        uint64  `json:"id"`
	FirstName string  `json:"fist_name"`
	LastName  *string `json:"last_name"`
}
