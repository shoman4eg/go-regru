package regru_test

import (
	"context"
	"errors"
	"fmt"

	regru "github.com/shoman4eg/go-regru"
)

// The public test account test/test works against the real API.
// No "Output:" comment, so `go test` compiles this example but does not call the network.
func Example() {
	client := regru.NewClient(nil).WithAuthPassword("test", "test")
	client.Lang = "ru"

	answer, _, err := client.Nop(context.Background())
	var apiErr *regru.ErrorResponse
	if errors.As(err, &apiErr) {
		fmt.Println("API error:", apiErr.Code, apiErr.Text)
		return
	}
	if err != nil {
		fmt.Println("transport error:", err)
		return
	}
	fmt.Println(answer.Login)
}
