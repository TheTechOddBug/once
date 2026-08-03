package command

import (
	"context"
	"fmt"

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
	host := args[0]

	app := ns.ApplicationByHost(host)
	if app == nil {
		return fmt.Errorf("no application found at host %q", host)
	}

	if !app.Running {
		return docker.ErrApplicationNotRunning
	}

	if err := ns.Setup(ctx); err != nil {
		return fmt.Errorf("%w: %w", docker.ErrSetupFailed, err)
	}

	settings := app.Settings
	if cmd.Flags().Changed("secret-key-base") {
		settings.Keys.SecretKeyBase = k.secretKeyBase
	}
	if cmd.Flags().Changed("vapid") {
		if err := settings.Keys.SetVAPIDKey(k.vapid); err != nil {
			return err
		}
	}

	oldSettings := app.Settings
	app.Settings = settings

	return runWithProgress("Setting keys for "+host, func(progress docker.DeployProgressCallback) error {
		if err := app.Deploy(ctx, progress); err != nil {
			app.Settings = oldSettings
			return fmt.Errorf("%w: %w", docker.ErrDeployFailed, err)
		}
		return nil
	})
}
