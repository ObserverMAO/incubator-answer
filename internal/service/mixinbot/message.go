package mixinbot

import "time"

const (
	maxRetries = 1 << 3
	baseDelay  = time.Second
)

const (
	MaxCardTitleLength   = 28
	MaxCardContentLength = 2 << 5
	PrefixTitle          = "💬"
	Ellipsis             = "..."
)
