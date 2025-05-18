package vercelkit

import "net/http"

type VercelHandler func(http.ResponseWriter, *http.Request)
