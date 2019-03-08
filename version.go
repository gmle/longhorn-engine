package main

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/rancher/longhorn-engine/controller/client"
	"github.com/rancher/longhorn-engine/meta"
)

func VersionCmd() cli.Command {
	return cli.Command{
		Name: "version",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name: "client-only",
			},
		},
		Action: func(c *cli.Context) {
			if err := version(c); err != nil {
				logrus.Fatalln("Error running info command:", err)
			}
		},
	}
}

type VersionOutput struct {
	ClientVersion *meta.VersionOutput `json:"clientVersion"`
	ServerVersion *meta.VersionOutput `json:"serverVersion"`
}

func version(c *cli.Context) error {
	clientVersion := meta.GetVersion()
	v := VersionOutput{ClientVersion: &clientVersion}

	if !c.Bool("client-only") {
		url := c.GlobalString("url")
		controllerClient := client.NewControllerClient(url)
		resp, err := controllerClient.GetVersion()
		if err != nil {
			return err
		}
		v.ServerVersion = &resp.VersionOutput
	}
	output, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}

	fmt.Println(string(output))
	return nil
}