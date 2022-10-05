/*
Copyright © 2022 Guilhem Lettron <guilhem@barpilot.io>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/go-logr/stdr"
	"github.com/guilhem/csi-runtime/driver"
	"github.com/spf13/cobra"

	"github.com/guilhem/csi-driver-git/pkg/identity"
	"github.com/guilhem/csi-driver-git/pkg/node"
)

// nodeCmd represents the node command
var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE:         nodeRun,
	SilenceUsage: true,
}

const driverName = "git.csi.barpilot.io"

var (
	endpoint          string
	nodeID            string
	maxVolumesPerNode int64
)

func init() {
	rootCmd.AddCommand(nodeCmd)

	nodeCmd.Flags().StringVar(&endpoint, "endpoint", "unix:///tmp/csi.sock", "CSI endpoint")

	nodeCmd.Flags().StringVar(&nodeID, "nodeid", "", "node ID")
	nodeCmd.MarkFlagRequired("nodeid")

	nodeCmd.Flags().Int64Var(&maxVolumesPerNode, "maxvolumespernode", 0, "limit of volumes per node")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// nodeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func nodeRun(cmd *cobra.Command, args []string) error {
	stdr.SetVerbosity(1)
	log := stdr.New(log.New(os.Stdout, "", log.Lshortfile))

	// react on CTRL+C and Sigterm
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return errors.New("no build info")
	}

	// log.Print(bi.Main.Version)

	// identity.New(driverName, version string, controller bool, manager identity.Interface)

	idm := identity.Identity{}

	// cm := controller.Controller{}

	ns, err := node.New(nodeID, maxVolumesPerNode)
	if err != nil {
		return err
	}

	driver, err := driver.New(ctx, endpoint, driverName, bi.Main.Version, idm, nil, ns, log)
	if err != nil {
		return err
	}

	driver.Start(ctx)

	return nil
}
