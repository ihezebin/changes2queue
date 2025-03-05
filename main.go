package main

import (
	"context"

	"github.com/ihezebin/changes2queue/cmd"
	"github.com/ihezebin/soup/logger"
)

func main() {
	ctx := context.Background()
	if err := cmd.Run(ctx); err != nil {
		logger.Fatalf(ctx, "cmd run error: %v", err)
	}
}
