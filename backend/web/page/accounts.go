package page

import (
	_ "embed"
	"html/template"
	"io"

	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/web/layout"
)

//go:embed accounts.html
var AccountsHTML string

type AccountsData struct {
	layout.BaseData

	Accounts []*model.Account
}

type AccountsPage struct {
	tmpl *template.Template
}

func NewAccounts() *AccountsPage {
	sources := []string{
		layout.BaseHTML,
		AccountsHTML,
	}

	tmpl := newTemplate("default", sources)
	page := AccountsPage{
		tmpl: tmpl,
	}
	return &page
}

func (p *AccountsPage) Render(w io.Writer, data AccountsData) error {
	return p.tmpl.ExecuteTemplate(w, "default", data)
}

func (p *AccountsPage) RenderAccounts(w io.Writer, data AccountsData) error {
	return p.tmpl.ExecuteTemplate(w, "accounts", data)
}
