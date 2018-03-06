//go:generate mockgen -package=mocks -destination=./mocks/registry.go github.com/SimonRichardson/alchemy/pkg/cluster/registry Registry

package registry

import (
	"github.com/SimonRichardson/alchemy/pkg/cluster/members"
)

type Registry interface {
	Add([]members.Member) error
	Remove([]members.Member) error
	Update([]members.Member) error
}
