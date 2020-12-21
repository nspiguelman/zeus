package responses

// EmptyResponse se usa para devolver respuestas con status 201
type EmptyResponse struct {
	Success bool
	Message string
}
