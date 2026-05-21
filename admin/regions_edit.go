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

type regionEditData struct {
	Country	*models.Country
	Region  *models.Region
	Error   string
	Success	string
}

func loadRegionsEdit(w http.ResponseWriter, db *sql.DB, r *http.Request) (regionEditData, error) {
	var data regionEditData
 
	countryID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid country id", http.StatusBadRequest)
		return data, err
	}

	regionID, err := uuid.Parse(r.PathValue("region_id"))
	if err != nil {
		http.Error(w, "invalid region id", http.StatusBadRequest)
		return data, err
	}
 
	data.Country, err = services.CountryReadByID(db, countryID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return data, err
	}
 
	data.Region, err = services.RegionReadByID(db, regionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return data, err
	}
 
	return data, nil
}

func regionsEdit(cfg *config.Config, db *sql.DB, tmpl *template.Template) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
		data, err := loadRegionsEdit(w, db, r)
		if err != nil {
			return
		}

        renderPage(w, tmpl, "regions_edit", data)
    }
}

func regionsEditPost(cfg *config.Config, db *sql.DB, tmpl *template.Template) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
		data, err := loadRegionsEdit(w, db, r)
		if err != nil {
			return
		}

		data.Region.Name      = r.FormValue("name")
		data.Region.NameAlt   = r.FormValue("name_alt")
		data.Region.IsEU      = r.FormValue("is_eu") == "on"
		data.Region.VATRate   = parseFloat(r.FormValue("vat_rate"))
		data.Region.IsEnabled = r.FormValue("is_enabled") == "on"

		if err := services.RegionUpdate(db, data.Region); err != nil {
			data.Error = friendlyError(err)
		} else {
			data.Success = "Region saved successfully."
		}
        
		renderPage(w, tmpl, "regions_edit", data)
    }
}

func regionsDeletePost(cfg *config.Config, db *sql.DB, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := loadRegionsEdit(w, db, r)
		if err != nil {
			return
		}

		if err := services.RegionDelete(db, data.Region.ID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/countries/"+data.Country.ID.String(), http.StatusSeeOther)
	}
}
