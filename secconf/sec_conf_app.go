package secconf

import (
	"errors"
	"fmt"
	"github.com/urfave/cli"
	"os"
	"sort"
)

func gen_check(c *cli.Context) error {
	output := c.String("output")
	file, err := createFile(output)
	if err != nil {
		return err
	}
	if file != nil {
		defer file.Close()
	}

	if !GenGeneratedAllCheck(file) {
		return errors.New("failed to generate sec config")
	}
	return nil
}

func gen_config(c *cli.Context) error {
	output := c.String("output")
	file, err := createFile(output)
	if err != nil {
		return err
	}
	if file != nil {
		defer file.Close()
	}

	if !GenGeneratedAllCommands(file) {
		return errors.New("failed to generate sec config")
	}
	return nil
}

func createFile(output string) (*os.File, error) {
	if output == "" {
		return nil, nil
	}

	file, error := os.Create(output)
	if error != nil {
		return nil, error
	}
	error = os.Chmod(output, 0700)
	if error != nil {
		file.Close()
		fmt.Println("Chmod err:", error)
		return nil, error
	}

	return file, nil
}

func newGenConfigCmd() *cli.Command {
	return &cli.Command{
		Name:   "gen_config",
		Usage:  "",
		Action: gen_config,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "output, o",
				Usage: "Writes the generated static baseline data to a specified file.",
			},
		},
	}
}

func newCenCheckCmd() *cli.Command {
	return &cli.Command{
		Name:   "gen_check",
		Usage:  "",
		Action: gen_check,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "output, o",
				Usage: "Writes the generated static baseline data to a specified file.",
			},
		},
	}
}

func newSecConfCmd() []cli.Command {
	return []cli.Command{
		*newGenConfigCmd(),
		*newCenCheckCmd(),
	}
}

func NewSecConfApp() *cli.App {
	app := cli.NewApp()
	app.Name = "sec_conf"
	app.Usage = "Security configuration tool"

	app.Commands = newSecConfCmd()
	sort.Sort(cli.CommandsByName(app.Commands))
	return app
}
