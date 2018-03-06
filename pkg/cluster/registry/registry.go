//go:generate mockgen -package=mocks -destination=./mocks/registry.go github.com/SimonRichardson/alchemy/pkg/cluster/registry Registry

package registry

type Registry interface{}
