package cmds

import (
	"io"

	"github.com/appscode/go/log"
	v "github.com/appscode/go/version"
	"github.com/spf13/cobra"
	"kube.ci/git-apiserver/pkg/cmds/server"
)

func NewCmdRun(out, errOut io.Writer, stopCh <-chan struct{}) *cobra.Command {
	o := server.NewStashOptions(out, errOut)

	cmd := &cobra.Command{
		Use:               "run",
		Short:             "Launch Stash Controller",
		Long:              "Launch Stash Controller",
		DisableAutoGenTag: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Infof("Starting operator version %s+%s ...", v.Version.Version, v.Version.CommitHash)

			if err := o.Complete(); err != nil {
				return err
			}
			if err := o.Validate(args); err != nil {
				return err
			}
			if err := o.Run(stopCh); err != nil {
				return err
			}
			return nil
		},
	}

	o.AddFlags(cmd.Flags())

	return cmd
}
