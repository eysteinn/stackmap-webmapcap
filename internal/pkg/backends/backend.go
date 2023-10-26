package backend

import "gitlab.com/EysteinnSig/stackmap-mapserver/internal/pkg/types"

/*type ProductModified struct {
	Product   string `yaml:"product"`
	Timestamp string `yaml:"timestamp"`
	Project   string `yaml:"project"`
}*/

type IBackend interface {
	OnReceived(*types.ProductModified) error
}
