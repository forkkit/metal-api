package service

// emptyBody is useful because with go-restful you cannot define an insert / update endpoint
// without specifying a payload for reading. it would immediately intercept the request and
// return 406: Not Acceptable to the client.
type EmptyBody struct{}
