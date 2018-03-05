package hashring

// Checksum computes a checksum for an instance of a HashRing. The
// checksum can be used to compare two rings for equality.
type Checksum interface {

	// Checksum calculates the checksum for the hashring that is passed in.
	Checksum(*HashRing) uint32
}
