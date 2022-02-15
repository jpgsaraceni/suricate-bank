package schema

type CachedResponse struct {
	Key            string
	ResponseStatus int
	ResponseBody   []byte
	EncodedRequest string
}
