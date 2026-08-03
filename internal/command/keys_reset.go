package command

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/basecamp/once/internal/docker"
)

type keysResetCommand struct {
	cmd           *cobra.Command
	secretKeyBase bool
	vapid         bool
}

func newKeysResetCommand() *keysResetCommand {
	k := &keysResetCommand{}
	k.cmd = &cobra.Command{
		Use:   "reset <host>",
		Short: "Regenerate secret keys and redeploy the application",
		Args:  cobra.ExactArgs(1),
		RunE:  WithNamespace(k.run),
	}

	k.cmd.Flags().BoolVar(&k.secretKeyBase, "secret-key-base", false, "regenerate the secret key base")
	k.cmd.Flags().BoolVar(&k.vapid, "vapid", false, "regenerate the VAPID key pair")
	k.cmd.MarkFlagsOneRequired("secret-key-base", "vapid")

	return k
}

// Private

func (k *keysResetCommand) run(ctx context.Context, ns *docker.Namespace, cmd *cobra.Command, args []string) error {
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
	if err := settings.Keys.Regenerate(k.secretKeyBase, k.vapid); err != nil {
		return fmt.Errorf("regenerating keys: %w", err)
	}

	oldSettings := app.Settings
	app.Settings = settings

	return runWithProgress("Resetting keys for "+host, func(progress docker.DeployProgressCallback) error {
		if err := app.Deploy(ctx, progress); err != nil {
			app.Settings = oldSettings
			return fmt.Errorf("%w: %w", docker.ErrDeployFailed, err)
		}
		return nil
	})
}
