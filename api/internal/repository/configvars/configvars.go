package configvars

import cvDom "github.com/TheAlpha16/isolet/api/internal/domain/configvars"

type cvRepo struct{}

func New() cvDom.Repository {
	return &cvRepo{}
}
