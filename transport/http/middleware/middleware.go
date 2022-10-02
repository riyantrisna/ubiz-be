package middleware

// Middleware is a struct that holds a middleware list
type Middleware struct {
	AuthMiddleware AuthMiddleware
	// add more middleware here
}
