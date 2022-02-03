package main

import (
	"os"
	"sync"

	"github.com/AkinoKaede/sstreamcry/common/net"
	"github.com/AkinoKaede/sstreamcry/shadowsocks"
	"github.com/urfave/cli/v2"
)

var wg sync.WaitGroup

func main() {
	app := &cli.App{
		Name: "sstreamcty",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "host",
				Aliases:  []string{"h"},
				Required: true,
			},
			&cli.IntFlag{
				Name:     "port",
				Aliases:  []string{"p"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "password",
				Aliases:  []string{"p"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "method",
				Aliases:  []string{"m"},
				Required: true,
			},
			&cli.IntFlag{
				Name:     "times",
				Aliases:  []string{"t"},
				Required: true,
			},
			&cli.IntFlag{
				Name: "threads",
			},
		},
		Action: func(c *cli.Context) error {
			account := shadowsocks.CreateAccount(c.String("password"), c.String("method"))
			dest := net.TCPDestination(net.ParseAddress(c.String("host")), net.Port(c.Int("port")))
			times := c.Int("times")
			if times < 1 {
				times = 1
			}

			threads := c.Int("threads")
			if threads < 1 {
				threads = 1
			}

			for i := 0; i < threads; i++ {
				wg.Add(1)

				go func() {
					shadowsocks.Boom(dest, account, times)
					wg.Done()
				}()
			}

			wg.Wait()

			return nil
		},
	}

	app.Run(os.Args)
}
