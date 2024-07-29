package utils

type subkey string
type ginContextKey string

const (
	SubKey        = subkey("sub")
	GinContextKey = ginContextKey("GinContextKey")
)
