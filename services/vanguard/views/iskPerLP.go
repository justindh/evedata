package views

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/antihax/evedata/services/vanguard"
	"github.com/antihax/evedata/services/vanguard/models"
	"github.com/antihax/evedata/services/vanguard/templates"
)

func init() {
	vanguard.AddRoute("iskPerLP", "GET", "/iskPerLP", iskPerLPPage)
	vanguard.AddRoute("iskPerLPCorpss", "GET", "/J/iskPerLPCorps", iskPerLPCorps)
	vanguard.AddRoute("iskPerLP", "GET", "/J/iskPerLP", iskPerLP)
}

func iskPerLPPage(w http.ResponseWriter, r *http.Request) {
	setCache(w, 60*60)
	p := newPage(r, "ISK Per Loyalty Point")

	templates.Templates = template.Must(template.ParseFiles("templates/iskPerLP.html", templates.LayoutPath))
	err := templates.Templates.ExecuteTemplate(w, "base", p)
	if err != nil {
		httpErr(w, err)
		return
	}
}

func iskPerLPCorps(w http.ResponseWriter, r *http.Request) {
	setCache(w, 60*60)
	v, err := models.GetISKPerLPCorporations()
	if err != nil {
		httpErr(w, err)
		return
	}

	json.NewEncoder(w).Encode(v)
}

func iskPerLP(w http.ResponseWriter, r *http.Request) {
	setCache(w, 60*30)
	q := r.FormValue("corp")
	v, err := models.GetISKPerLP(q)
	if err != nil {
		httpErr(w, err)
		return
	}

	json.NewEncoder(w).Encode(v)
}
