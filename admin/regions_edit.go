package admin

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"go.etcd.io/bbolt"

	"gocart/config"
	"gocart/models"
)

type regionEditData struct {
    Country *models.Country
    Region  *models.Region
    Index   int
    Error   string
    Success string
}

func loadRegionsEdit(w http.ResponseWriter, db *bbolt.DB, r *http.Request) (regionEditData, uuid.UUID, error) {
	var data regionEditData

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return data, uuid.Nil, err
	}

	data.Index, err = strconv.Atoi(r.PathValue("index"))
	if err != nil {
		http.Error(w, "invalid index", http.StatusBadRequest)
		return data, uuid.Nil, err
	}

	data.Country, err = models.CountryReadByID(db, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return data, uuid.Nil, err
	}

	if data.Index < 0 || data.Index >= len(data.Country.Regions) {
		http.Error(w, "region not found", http.StatusNotFound)
		return data, uuid.Nil, fmt.Errorf("region not found")
	}

	data.Region = data.Country.Regions[data.Index]

	return data, id, nil
}

func regionsEdit(cfg *config.Config, db *bbolt.DB, tmpl *template.Template) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
		data, _, err := loadRegionsEdit(w, db, r)
		if err != nil {
			return
		}

        renderPage(w, tmpl, "regions_edit", data)
    }
}

func regionsEditPost(cfg *config.Config, db *bbolt.DB, tmpl *template.Template) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
		data, _, err := loadRegionsEdit(w, db, r)
		if err != nil {
			return
		}

		data.Region.Name      = r.FormValue("name")
		data.Region.NameAlt   = r.FormValue("name_alt")
		data.Region.IsEU      = r.FormValue("is_eu") == "on"
		data.Region.VATRate   = parseFloat(r.FormValue("vat_rate"))
		data.Region.IsEnabled = r.FormValue("is_enabled") == "on"

		if err := models.CountryUpdate(db, data.Country); err != nil {
			data.Error = friendlyError(err)
		} else {
			data.Success = "Region saved successfully."
		}
        
		renderPage(w, tmpl, "regions_edit", data)
    }
}

func regionsDeletePost(cfg *config.Config, db *bbolt.DB, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, id, err := loadRegionsEdit(w, db, r)
		if err != nil {
			return
		}

		data.Country.Regions = append(data.Country.Regions[:data.Index], data.Country.Regions[data.Index+1:]...)
		if err := models.CountryUpdate(db, data.Country); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/countries/"+id.String(), http.StatusSeeOther)
	}
}
