package admin

import (
    "net/http"
    "html/template"
    "strconv"

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

type regionEditData struct {
    Country *models.Country
    Region  *models.Region
    Index   int
    Error   string
    Success string
}

func getRegionNew(cfg *config.Config, db *bbolt.DB, tmpl *template.Template) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        id, err := uuid.Parse(r.PathValue("id"))

        if err != nil {
            http.Error(w, "invalid id", http.StatusBadRequest)
            return
        }

        c, err := models.CountryReadByID(db, id)
        if err != nil {
            http.Error(w, err.Error(), http.StatusNotFound)
            return
        }

        render(w, tmpl, "regions-new", regionNewData{
            Country: c,
            Region:  &models.Region{},
        })
    }
}

func postRegionNew(cfg *config.Config, db *bbolt.DB, tmpl *template.Template) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        id, err := uuid.Parse(r.PathValue("id"))
        if err != nil {
            http.Error(w, "invalid id", http.StatusBadRequest)
            return
        }

        c, err := models.CountryReadByID(db, id)
        if err != nil {
            http.Error(w, err.Error(), http.StatusNotFound)
            return
        }

        region := &models.Region{
            Name:      r.FormValue("name"),
            NameAlt:   r.FormValue("name_alt"),
            IsEU:      r.FormValue("is_eu") == "on",
            VATRate:   parseFloat(r.FormValue("vat_rate")),
            IsEnabled: r.FormValue("is_enabled") == "on",
        }

        c.Regions = append(c.Regions, region)
        if err := models.CountryUpdate(db, c); err != nil {
            render(w, tmpl, "regions-new", regionNewData{
                Country: c,
                Region:  region,
                Error:   friendlyError(err),
            })

            return
        }

        http.Redirect(w, r, "/countries/" + id.String(), http.StatusSeeOther)
    }
}

func getRegionEdit(cfg *config.Config, db *bbolt.DB, tmpl *template.Template) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        id, err := uuid.Parse(r.PathValue("id"))
        if err != nil {
            http.Error(w, "invalid id", http.StatusBadRequest)
            return
        }

        index, err := strconv.Atoi(r.PathValue("index"))
        if err != nil {
            http.Error(w, "invalid index", http.StatusBadRequest)
            return
        }

        c, err := models.CountryReadByID(db, id)
        if err != nil {
            http.Error(w, err.Error(), http.StatusNotFound)
            return
        }

        if index < 0 || index >= len(c.Regions) {
            http.Error(w, "region not found", http.StatusNotFound)
            return
        }

        render(w, tmpl, "regions-edit", regionEditData{
            Country: c,
            Region:  c.Regions[index],
            Index:   index,
        })
    }
}

func postRegionEdit(cfg *config.Config, db *bbolt.DB, tmpl *template.Template) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        id, err := uuid.Parse(r.PathValue("id"))
        if err != nil {
            http.Error(w, "invalid id", http.StatusBadRequest)
            return
        }

        index, err := strconv.Atoi(r.PathValue("index"))
        if err != nil {
            http.Error(w, "invalid index", http.StatusBadRequest)
            return
        }

        c, err := models.CountryReadByID(db, id)
        if err != nil {
            http.Error(w, err.Error(), http.StatusNotFound)
            return
        }

        if index < 0 || index >= len(c.Regions) {
            http.Error(w, "region not found", http.StatusNotFound)
            return
        }

        c.Regions[index].Name      = r.FormValue("name")
        c.Regions[index].NameAlt   = r.FormValue("name_alt")
        c.Regions[index].IsEU      = r.FormValue("is_eu") == "on"
        c.Regions[index].VATRate   = parseFloat(r.FormValue("vat_rate"))
        c.Regions[index].IsEnabled = r.FormValue("is_enabled") == "on"

        if err := models.CountryUpdate(db, c); err != nil {
            render(w, tmpl, "regions-edit", regionEditData{
                Country: c,
                Region:  c.Regions[index],
                Index:   index,
                Error:   friendlyError(err),
            })

            return
        }

        http.Redirect(w, r, "/countries/" + id.String(), http.StatusSeeOther)
    }
}

func postRegionDelete(cfg *config.Config, db *bbolt.DB, tmpl *template.Template) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        id, err := uuid.Parse(r.PathValue("id"))
        if err != nil {
            http.Error(w, "invalid id", http.StatusBadRequest)
            return
        }

        index, err := strconv.Atoi(r.PathValue("index"))
        if err != nil {
            http.Error(w, "invalid index", http.StatusBadRequest)
            return
        }

        c, err := models.CountryReadByID(db, id)
        if err != nil {
            http.Error(w, err.Error(), http.StatusNotFound)
            return
        }

        if index < 0 || index >= len(c.Regions) {
            http.Error(w, "region not found", http.StatusNotFound)
            return
        }

        c.Regions = append(c.Regions[:index], c.Regions[index+1:]...)
        if err := models.CountryUpdate(db, c); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
		
        http.Redirect(w, r, "/countries/" + id.String(), http.StatusSeeOther)
    }
}
