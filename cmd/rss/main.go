package main

import (
	"context"
	"fmt"

	"github.com/moviegeek/ptrss"
)

func main() {
	ctx := context.Background()

	err := ptrss.UpdateRss(ctx, ptrss.PubSubMessage{})
	if err != nil {
		fmt.Println(err)
	}
}
