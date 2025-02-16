package main

import (
	"chatgpt-adapter/wire"
	"github.com/iocgo/sdk"
	"github.com/iocgo/sdk/errors"
)

//go:generate goversioninfo -arm -o=resource_windows.syso -icon=favicon.ico

func main() {
	ctx := errors.New(nil)
	{
		if err := errors.Try1(ctx, func() (c *sdk.Container, err error) {
			c = sdk.NewContainer()
			err = wire.Injects(c)
			return
		}).Run(); err != nil {
			panic(err)
		}
	}
}
