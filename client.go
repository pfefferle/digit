package digit

import (
	"github.com/benpate/derp"
	"github.com/benpate/remote"
)

// Lookup resolves a WebFinger account (a URL, an email-style handle, or a
// Fediverse address) into its Resource, querying each candidate endpoint until
// one succeeds. It returns an error if no endpoint can be reached.
func Lookup(url string, options ...remote.Option) (Resource, error) {

	webFingerServerURLs := ParseAccount(url)
	result := NewResource(url)

	for _, webFingerServerURL := range webFingerServerURLs {

		txn := remote.Get(webFingerServerURL).
			With(options...).
			Result(&result)

		if err := txn.Send(); err == nil {
			return result, nil
		}
	}

	return result, derp.Internal("digit.Lookup", "Unable to load resource", url, webFingerServerURLs)
}
