package main

import (
	"log"
	"os"
	"sync"

	"github.com/AkinoKaede/sstreamcry/common/net"
	"github.com/AkinoKaede/sstreamcry/shadowsocks"
	"github.com/urfave/cli/v2"
)

var wg sync.WaitGroup

func main() {
	app := &cli.App{
		Name: "sstreamcry",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "host",
				Required: true,
			},
			&cli.IntFlag{
				Name:     "port",
				Aliases:  []string{"p"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "password",
				Aliases:  []string{"pwd"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "method",
				Aliases:  []string{"m"},
				Required: true,
			},
			&cli.IntFlag{
				Name:        "times",
				Aliases:     []string{"t"},
				DefaultText: "1",
				Required:    true,
			},
			&cli.IntFlag{
				Name:        "threads",
				Aliases:     []string{"tr"},
				DefaultText: "1",
			},
		},
		Action: func(c *cli.Context) error {
			account, err := shadowsocks.CreateAccount(c.String("password"), c.String("method"))
			if err != nil {
				return err
			}

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
					err := shadowsocks.Boom(dest, *account, times)
					log.Println(err)
					wg.Done()
				}()
			}

			wg.Wait()

			return nil
		},
		ExitErrHandler: func(_ *cli.Context, err error) {
			if err != nil {
				log.Println(err)
			}
		},
	}

	app.Run(os.Args)
}
