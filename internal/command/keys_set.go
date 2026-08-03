package command

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/basecamp/once/internal/docker"
)

type keysSetCommand struct {
	cmd           *cobra.Command
	secretKeyBase string
	vapid         string
}

func newKeysSetCommand() *keysSetCommand {
	k := &keysSetCommand{}
	k.cmd = &cobra.Command{
		Use:   "set <host>",
		Short: "Set secret keys and redeploy the application",
		Args:  cobra.ExactArgs(1),
		RunE:  WithNamespace(k.run),
	}

	k.cmd.Flags().StringVar(&k.secretKeyBase, "secret-key-base", "", "new secret key base")
	k.cmd.Flags().StringVar(&k.vapid, "vapid", "", "new VAPID private key (the public key is derived from it)")
	k.cmd.MarkFlagsOneRequired("secret-key-base", "vapid")

	return k
}

// Private

func (k *keysSetCommand) run(ctx context.Context, ns *docker.Namespace, cmd *cobra.Command, args []string) error {
	return changeKeys(ctx, ns, args[0], "Setting keys for", func(keys *docker.Keys) error {
		if cmd.Flags().Changed("secret-key-base") {
			keys.SecretKeyBase = k.secretKeyBase
		}
		if cmd.Flags().Changed("vapid") {
			return keys.SetVAPIDKey(k.vapid)
		}
		return nil
	})
}
