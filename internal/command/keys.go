package command

import "github.com/spf13/cobra"

type keysCommand struct {
	cmd *cobra.Command
}

func newKeysCommand() *keysCommand {
	k := &keysCommand{}
	k.cmd = &cobra.Command{
		Use:   "keys",
		Short: "Manage application secret keys",
	}

	k.cmd.AddCommand(newKeysResetCommand().cmd)
	k.cmd.AddCommand(newKeysSetCommand().cmd)

	return k
}
