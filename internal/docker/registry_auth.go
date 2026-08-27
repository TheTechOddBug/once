package docker

import (
	"encoding/base64"
	"encoding/json"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
)

// RegistryHost returns the registry host an image ref points at,
// such as "ghcr.io" or "index.docker.io".
func RegistryHost(imageRef string) (string, error) {
	ref, err := name.ParseReference(imageRef)
	if err != nil {
		return "", err
	}
	return ref.Context().RegistryStr(), nil
}

// registryAuthFor returns a base64-encoded JSON auth string for pulling the
// application's image, suitable for use in ImagePullOptions.RegistryAuth.
// Credentials from the application settings are used when they are scoped to
// the image's registry host; otherwise the docker keychain is consulted.
// Returns "" on any error or missing credentials, falling back to anonymous access.
func registryAuthFor(settings ApplicationSettings) string {
	if auth := credentialAuth(settings.Registry, settings.Image); auth != "" {
		return auth
	}
	return keychainAuth(settings.Image)
}

// Helpers

func credentialAuth(registry RegistrySettings, imageName string) string {
	host, err := RegistryHost(imageName)
	if err != nil || registry.Empty() || registry.Host != host {
		return ""
	}
	return encodeAuthConfig(authn.AuthConfig{
		Username: registry.Username,
		Password: registry.Password,
	})
}

func keychainAuth(imageName string) string {
	ref, err := name.ParseReference(imageName)
	if err != nil {
		return ""
	}
	authenticator, err := authn.DefaultKeychain.Resolve(ref.Context())
	if err != nil || authenticator == authn.Anonymous {
		return ""
	}
	cfg, err := authenticator.Authorization()
	if err != nil {
		return ""
	}
	return encodeAuthConfig(*cfg)
}

func encodeAuthConfig(cfg authn.AuthConfig) string {
	data, err := json.Marshal(cfg)
	if err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(data)
}
