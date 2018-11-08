package cmds

import (
	"io"

	"github.com/appscode/go/log"
	v "github.com/appscode/go/version"
	"github.com/kube-ci/git-apiserver/pkg/cmds/server"
	"github.com/spf13/cobra"
)

func NewCmdRun(out, errOut io.Writer, stopCh <-chan struct{}) *cobra.Command {
	o := server.NewGitAPIServerOptions(out, errOut)

	cmd := &cobra.Command{
		Use:               "run",
		Short:             "Launch git apiserver",
		Long:              "Launch git apiserver",
		DisableAutoGenTag: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Infof("Starting server version %s+%s ...", v.Version.Version, v.Version.CommitHash)

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
