package main

// customized contextkey type to avoid naming collisions with third parties.

type contextKey string

const isAuthenticatedContextKey = contextKey("isAuthenticated")
