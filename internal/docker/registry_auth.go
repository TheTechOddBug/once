package docker

import (
	"encoding/base64"
	"encoding/json"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
)

// registryAuthFor returns a base64-encoded JSON auth string for pulling the
// application's image, suitable for use in ImagePullOptions.RegistryAuth.
// Credentials from the application settings are used when they are scoped to
// the image's repository; otherwise the docker keychain is consulted.
// Returns "" on any error or missing credentials, falling back to anonymous access.
func registryAuthFor(settings ApplicationSettings) string {
	if auth := credentialAuth(settings.Registry, settings.Image); auth != "" {
		return auth
	}
	return keychainAuth(settings.Image)
}

// Helpers

func credentialAuth(registry RegistrySettings, imageName string) string {
	if registry.Empty() || !sameRepository(registry.Image, imageName) {
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

// sameRepository reports whether two image refs point at the same repository,
// ignoring tags and digests. Unparsable refs never match.
func sameRepository(a, b string) bool {
	refA, err := name.ParseReference(a)
	if err != nil {
		return false
	}
	refB, err := name.ParseReference(b)
	if err != nil {
		return false
	}
	return refA.Context().Name() == refB.Context().Name()
}

func encodeAuthConfig(cfg authn.AuthConfig) string {
	data, err := json.Marshal(cfg)
	if err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(data)
}
