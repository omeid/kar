// +build kar

package main

import (
	"github.com/omeid/gonzo/context"
	"github.com/omeid/kar"
	"github.com/omeid/kargar"
)

func init() {

	kar.Run(func(build *kargar.Build) error {

		return build.Add(
			kargar.Task{

				Name:  "say-hello",
				Usage: "This tasks is self-documented, it says hello for every second.",

				Action: func(ctx context.Context) error {
					ctx.Info("Hello!")
					return nil
				},
			})
	})
}
