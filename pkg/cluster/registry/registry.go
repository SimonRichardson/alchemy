//go:generate mockgen -package=mocks -destination=./mocks/registry.go github.com/SimonRichardson/alchemy/pkg/cluster/registry Registry

package registry

// Address represents a url that can be used for a host and port inside the
// registry.
type Address string

// HostPort returns the address host and port (host:port)
func (a Address) HostPort() string {
	return string(a)
}

type Key interface {

	// Name returns the registry key
	Name() string

	// Type represents different types of key variants
	Type() string

	// Address defines the url of the key.
	Address() Address

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
}
