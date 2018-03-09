//go:generate mockgen -package=mocks -destination=./mocks/registry.go github.com/SimonRichardson/alchemy/pkg/cluster/registry Registry

package registry

type Key interface {

	// Name returns the registry key
	Name() string

	// Type represents different types of key variants
	Type() string

	// Address defines the url of the key.
	Address() string

	// Tags returns any associated tags of the key
	Tags() map[string]string
}

type Registry interface {

	// Add a key to the registry. Adding the key multiple times should not
	// change the underlying storage.
	// Returns true if the key was added to the registry
	Add(Key) bool

	// Remove a key from the registry. Removing a key multiple times will
	// not change the underlying storage.
	// Returns true if the key was remove from the registry
	Remove(Key) bool

	// Update updates a key in place. Updating the key should be done in place.
	// Returns true if the key was updated to the registry
	Update(Key) bool

	// Locate attempts to locate a key depending on the value of the string
	// passed in. The string doesn't represent the hash, but potentially a path
	// or address that will return the same key if the underlying storage is
	// stable.
	// Returns true if the key was located or not from the registry
	Locate(string) (Key, bool)

	// Info returns back the information for a particular key type
	// Returns an error if it can't gather the information about a specific key
	// type
	Info(string) (Info, error)

	// Checksum returns back a checksum for a particular key type
	// Returns an error if it can't gather the checksum about a specific key
	// type
	Checksum(string) (string, error)
}

// Info represents information for a registry key type
type Info struct {
	Checksum string
	Hashes   map[string]string
	Keys     map[string][]Key
}
