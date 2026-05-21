package admin

import (
	"database/sql"
	"html/template"
	"net/http"
 
	"github.com/google/uuid"

	"gocart/config"
	"gocart/models"
	"gocart/services"
)

type regionNewData struct {
    Country	*models.Country
    Region  *models.Region
    Error   string
}

func loadRegionsNew(w http.ResponseWriter, db *sql.DB, r *http.Request) (regionNewData, uuid.UUID, error) {
	var data regionNewData

	countryID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid country id", http.StatusBadRequest)
		return data, uuid.Nil, err
	}

	data.Country, err = services.CountryReadByID(db, countryID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return data, uuid.Nil, err
	}

    data.Region = &models.Region{}

	return data, countryID, nil
}

func regionsNew(cfg *config.Config, db *sql.DB, tmpl *template.Template) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
		data, _, err := loadRegionsNew(w, db, r)
		if err != nil {
			return
		}

        renderPage(w, tmpl, "regions_new", data)
    }
}

func regionsNewPost(cfg *config.Config, db *sql.DB, tmpl *template.Template) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
		data, countryID, err := loadRegionsNew(w, db, r)
		if err != nil {
			return
		}

		data.Region = &models.Region{
			CountryID: countryID,
			Name:      r.FormValue("name"),
			NameAlt:   r.FormValue("name_alt"),
			IsEU:      r.FormValue("is_eu") == "on",
			VATRate:   parseFloat(r.FormValue("vat_rate")),
		}
		data.Region.IsEnabled = true

		if err := services.RegionCreate(db, data.Region); err != nil {
			data.Error = friendlyError(err)
			renderPage(w, tmpl, "regions_new", data)
			return
		}

		http.Redirect(w, r, "/countries/"+countryID.String(), http.StatusSeeOther)
    }
}
