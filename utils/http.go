package utils

import (
	"net/http"
	"time"
)

var HTTPClient = &http.Client{Timeout: 60 * time.Second}
