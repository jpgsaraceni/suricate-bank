package schema

type CachedResponse struct {
	Key            string
	ResponseStatus int
	ResponseBody   []byte
}

func NewCachedResponse(key string, status int, body []byte) CachedResponse {
	return CachedResponse{
		Key:            key,
		ResponseStatus: status,
		ResponseBody:   body,
	}
}
