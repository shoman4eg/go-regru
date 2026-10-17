package regru

import (
	"context"
	"testing"
)

func TestHostingService(t *testing.T) {
	c, mux := setup(t)
	handle(t, mux, "/hosting/nop", nil, `{}`)
	if _, _, err := c.Hosting.Nop(context.Background()); err != nil {
		t.Fatal(err)
	}
}
