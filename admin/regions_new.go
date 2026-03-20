package admin

import (
    "html/template"
    "net/http"

    "github.com/google/uuid"
    "go.etcd.io/bbolt"

    "gocart/config"
    "gocart/models"
)

type regionNewData struct {
    Country *models.Country
    Region  *models.Region
    Error   string
}

func regionsNew(cfg *config.Config, db *bbolt.DB, tmpl *template.Template) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
		var data regionNewData

        id, err := uuid.Parse(r.PathValue("id"))
        if err != nil {
            http.Error(w, "invalid id", http.StatusBadRequest)
            return
        }

        data.Country, err = models.CountryReadByID(db, id)
        if err != nil {
            http.Error(w, err.Error(), http.StatusNotFound)
            return
        }

        renderPage(w, tmpl, "regions_new", data)
    }
}

func regionsNewPost(cfg *config.Config, db *bbolt.DB, tmpl *template.Template) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
		var data regionNewData

        id, err := uuid.Parse(r.PathValue("id"))
        if err != nil {
            http.Error(w, "invalid id", http.StatusBadRequest)
            return
        }

        data.Country, err = models.CountryReadByID(db, id)
        if err != nil {
            http.Error(w, err.Error(), http.StatusNotFound)
            return
        }

		data.Region = &models.Region{
			Name:      r.FormValue("name"),
			NameAlt:   r.FormValue("name_alt"),
			IsEU:      r.FormValue("is_eu") == "on",
			VATRate:   parseFloat(r.FormValue("vat_rate")),
			IsEnabled: r.FormValue("is_enabled") == "on",
		}

		data.Country.Regions = append(data.Country.Regions, data.Region)

        if err := models.CountryUpdate(db, data.Country); err != nil {
			data.Error = friendlyError(err)
			renderPage(w, tmpl, "regions_new", data)
		} else {
			http.Redirect(w, r, "/countries/" + id.String(), http.StatusSeeOther)
		}
    }
}
