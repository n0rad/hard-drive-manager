package checksum

import (
	"github.com/n0rad/hard-disk-manager/pkg/checksum/integrity"
	"github.com/n0rad/hard-disk-manager/pkg/config"
	"github.com/spf13/cobra"
)

func listCommand(conf *config.GlobalConfig) *cobra.Command {
	var reverse bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list files",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if err := runCmdForPath(conf, arg, func(pathConf config.PathConfig, d integrity.Directory) func(path string) error {
					if reverse {
						d.Exclusive = !d.Exclusive
					}
					return d.List
				}); err != nil {
					return err
				}
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&reverse, "reverse", "r", false, "Reverse regex match")
	return cmd
}
