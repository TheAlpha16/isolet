package middleware

var skipPaths = map[string]struct{}{
	"/ping": {},
}
