package repository

import (
	"github.com/itohio/phingo/pkg/types"
)

func (r *repository) SetConfig(cfg *types.Config) {
	r.config = cfg
}

func (r *repository) SetAccount(acc ...*types.Account) {
}

func (r *repository) SetClient(cl ...*types.Client) {
}

func (r *repository) SetProject(prj ...*types.Project) {
}

func (r *repository) SetInvoice(inv ...*types.Invoice) {
}

func (r *repository) SetTemplate(tpl ...*types.Template) {
}
