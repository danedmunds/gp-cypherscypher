package main

import (
	"fmt"
	"log"
	"os"
	"unicode"

	"github.com/urfave/cli"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func main() {
	app := cli.NewApp()
	app.Name = "ciphers"
	app.Usage = "cipher text"

	var shift int
	var keyword string
	var rails int
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "decipher, d",
			Usage: "Decipher the input",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "caesar",
			Aliases: []string{"c"},
			Usage:   "Caesar cipher the input",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:        "shift, s",
					Value:       13,
					Usage:       "Right shift by `SHIFT`",
					Destination: &shift,
				},
			},
			Action: func(c *cli.Context) error {
				doIt(c, Caesar(shift))
				return nil
			},
		},
		{
			Name:    "rot13",
			Aliases: []string{"r"},
			Usage:   "Rot13 cipher the input",
			Action: func(c *cli.Context) error {
				doIt(c, Rot13())
				return nil
			},
		},
		{
			Name:    "keyword",
			Aliases: []string{"k"},
			Usage:   "Keyword cipher the input",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "keyword, k",
					Usage:       "Cipher using `KEYWORD`",
					Destination: &keyword,
				},
			},
			Action: func(c *cli.Context) error {
				doIt(c, Keyword(keyword))
				return nil
			},
		},
		{
			Name:    "railfence",
			Aliases: []string{"r"},
			Usage:   "Rail Fence cipher the input",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:        "rails, r",
					Usage:       "Cipher using `RAILS` rails",
					Destination: &rails,
					Value:       3,
				},
			},
			Action: func(c *cli.Context) error {
				doIt(c, RailFence(rails))
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func doIt(ctx *cli.Context, cipher Cipher) {
	var inOut InOutFunc
	// TODO make sure about this, had it backwards before
	if ctx.GlobalBool("decipher") {
		inOut = cipher.Decipher
	} else {
		inOut = cipher.Encipher
	}

	t := transform.Chain(
		norm.NFKD,
		runes.Remove(runes.In(unicode.Mark)),
		runes.Map(func(r rune) rune {
			return unicode.ToUpper(r)
		}),
	)

	err := inOut(transform.NewReader(os.Stdin, t), os.Stdout)
	if err != nil {
		fmt.Printf("%+v\n", err)
		panic(err)
	}
}
