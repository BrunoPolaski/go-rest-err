# go-rest-err

Small helper library for REST-style errors in Go.

## Usage

```go
package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	resterr "github.com/BrunoPolaski/go-rest-err/rest_err"
)

type ErrorResponse struct {
	Error *resterr.RestErr `json:"error"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Imagine this came from your data layer
	dbErr := errors.New("db connection failed")

	// Wrap the root error
	err := resterr.NewInternalServerError("could not fetch user").WithCause(dbErr)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Code)
	_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err})
}

func main() {
	log.Fatal(http.ListenAndServe(":8080", http.HandlerFunc(handler)))
}
```

### Checking error chains

```go
if errors.Is(err, dbErr) {
	// handle root cause
}

var restErr *resterr.RestErr
if errors.As(err, &restErr) {
	// handle REST error
}
```
