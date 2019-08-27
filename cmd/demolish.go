// Copyright 2017 Palantir Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/palantir/bouncer/bouncer"
	"github.com/palantir/bouncer/demolish"
)

var demolishCmd = &cobra.Command{
	Use:   "demolish",
	Short: "Run bouncer in demolish",
	Long:  `Run bouncer in demolish mode, where we demolish all nodes from the list of ASGs and wait for the new instances to be ready.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(logLevelFromViper())

		log.Debug("demolish called")
		if log.GetLevel() == log.DebugLevel {
			cmd.DebugFlags()
			viper.Debug()
		}

		asgsString := viper.GetString("demolish.asgs")
		if asgsString == "" {
			log.Fatal("You must specify ASGs to demolish nodes from (in a comma-delimited list)")
		}

		commandString := viper.GetString("demolish.command")
		noop := viper.GetBool("demolish.noop")
		force := viper.GetBool("demolish.force")
		termHook := viper.GetString("terminate-hook")
		pendHook := viper.GetString("pending-hook")
		timeout := timeoutFromViper()

		log.Debugf("Binding vars, got %+v %+v %+v %+v", asgsString, noop, version, commandString)

		log.Info("Beginning bouncer demolish run")

		var defCap int64
		defCap = 1
		opts := bouncer.RunnerOpts{
			Noop:            noop,
			Force:           force,
			AsgString:       asgsString,
			CommandString:   commandString,
			DefaultCapacity: &defCap,
			TerminateHook:   termHook,
			PendingHook:     pendHook,
			ItemTimeout:     timeout,
		}

		r, err := demolish.NewRunner(&opts)
		if err != nil {
			log.Fatal(errors.Wrap(err, "error initializing runner"))
		}

		r.MustValidatePrereqs()

		err = r.Run()
		if err != nil {
			log.Fatal(errors.Wrap(err, "error in run"))
		}
	},
}

func init() {
	RootCmd.AddCommand(demolishCmd)

	demolishCmd.Flags().BoolP("noop", "n", false, "Run this in noop mode, and only print what you would do")
	err := viper.BindPFlag("demolish.noop", demolishCmd.Flags().Lookup("noop"))
	if err != nil {
		log.Fatal(errors.Wrap(err, "Binding PFlag 'noop' to viper var 'demolish.noop' failed: %s"))
	}

	demolishCmd.Flags().StringP("asgs", "a", "", "ASGs to check for nodes to cycle in")
	err = viper.BindPFlag("demolish.asgs", demolishCmd.Flags().Lookup("asgs"))
	if err != nil {
		log.Fatal(errors.Wrap(err, "Binding PFlag 'asgs' to viper var 'demolish.asgs' failed: %s"))
	}

	demolishCmd.Flags().StringP("preterminatecall", "p", "", "External command to run before host is removed from its ELB & terminate process begins")
	err = viper.BindPFlag("demolish.command", demolishCmd.Flags().Lookup("preterminatecall"))
	if err != nil {
		log.Fatal(errors.Wrap(err, "Binding PFlag 'command' to viper var 'demolish.command' failed: %s"))
	}

	demolishCmd.Flags().BoolP("force", "f", false, "Force all nodes to be recycled, even if they're running the latest launch config")
	err = viper.BindPFlag("demolish.force", demolishCmd.Flags().Lookup("force"))
	if err != nil {
		log.Fatal(errors.Wrap(err, "Binding PFlag 'force' to viper var 'demolish.force' failed: %s"))
	}
}
