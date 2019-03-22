package responses

type ResponseModel interface {
	MarshalToJSON()([]byte, error)
}
