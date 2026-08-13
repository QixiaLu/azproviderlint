package response

func WasNotFound(statusCode int) bool {
	return statusCode == 404
}
