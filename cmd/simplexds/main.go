package main

import (
	"context"

	"github.com/wozz/simplexds"
)

func main() {
	ctx := context.Background()
	simplexds.Run(ctx)
	<-ctx.Done()
}
