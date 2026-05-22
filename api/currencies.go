package api

import (
    "fmt"
    "encoding/json"
    "net/http"

    "github.com/google/uuid"
    "github.com/jmoiron/sqlx"

    "gocart/models"
    "gocart/utils"
)

func CreateCurrency(db *sqlx.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")

        var c_body models.Currency
        var c_created models.Currency

        if err := json.NewDecoder(r.Body).Decode(&c_body); err != nil {
            errRes := utils.ErrorResponse{}
            errRes.SetError("JSON format is invalid")
            errRes.Write(w, http.StatusBadRequest)
            return
        }


        if errRes := c_body.Validate(); errRes != nil {
            errRes.Write(w, http.StatusBadRequest)
            return
        }

        // We don't insert the ID and timestamps from the body
        // We let the DB generate default ones instead, and return what DB generated
        query := `INSERT INTO currencies (
            iso_code,
            flag_name,
            name_en,
            name_es,
            decimals,
            fx_rate,
            is_enabled
        ) VALUES (
            :iso_code,
            :flag_name,
            :name_en,
            :name_es,
            :decimals,
            :fx_rate,
            :is_enabled
        ) RETURNING *`

        stmt, _ := db.PrepareNamed(query)
        if err := stmt.Get(&c_created, c_body); err != nil {
            errRes := utils.ErrorResponse{}
            errRes.SetError("Failed to create item")
            errRes.Write(w, http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(c_created)
    }
}

func ReadCurrencyByID(db *sqlx.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")

        var c models.Currency
        
        id, err := uuid.Parse(r.PathValue("id"))
        if err != nil {
            errRes := utils.ErrorResponse{}
            errRes.SetError("ID format is invalid")
            errRes.Write(w, http.StatusBadRequest)
            return
        }

        query := `SELECT * FROM currencies WHERE id = $1`

        if err := db.Get(&c, query, id); err != nil {
            errRes := utils.ErrorResponse{}
            errRes.SetError("Item not found")
            errRes.Write(w, http.StatusNotFound)
            return
        }

        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(c)
    }
}

func ReadCurrencies(db *sqlx.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")

        currencies := []models.Currency{}

        params := utils.ParseQueryParams(r)

        params.Limit    = utils.SanitizeInt(params.Limit, 1, 500, 300)
        params.Offset   = utils.SanitizeInt(params.Offset, 0, 100000, 0)
        params.OrderBy  = utils.SanitizeString(params.OrderBy, []string{"iso_code", "name_en", "name_es", "is_enabled"}, "name_en")
        params.OrderDir = utils.SanitizeString(params.OrderDir, []string{"asc", "desc"}, "asc")

        query := fmt.Sprintf(`SELECT *
            FROM currencies
            ORDER BY %s %s
            LIMIT $1 OFFSET $2`, params.OrderBy, params.OrderDir)

        if err := db.Select(&currencies, query, params.Limit, params.Offset); err != nil {
            errRes := utils.ErrorResponse{}
            errRes.SetError("Failed to get items")
            errRes.Write(w, http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(currencies)
    }
}

func UpdateCurrencyByID(db *sqlx.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")

        var c_body models.Currency
        var c_updated models.Currency
        
        id, err := uuid.Parse(r.PathValue("id"))
        if err != nil {
            errRes := utils.ErrorResponse{}
            errRes.SetError("ID format is invalid")
            errRes.Write(w, http.StatusBadRequest)
            return
        }

        if err := json.NewDecoder(r.Body).Decode(&c_body); err != nil {
            errRes := utils.ErrorResponse{}
            errRes.SetError("JSON format is invalid")
            errRes.Write(w, http.StatusBadRequest)
            return
        }

        c_body.ID = id // Ignore any ID in the body and replace with the ID from the path

        if errRes := c_body.Validate(); errRes != nil {
            errRes.Write(w, http.StatusBadRequest)
            return
        }

        query := `UPDATE currencies SET
            iso_code   = :iso_code,
            flag_name  = :flag_name,
            name_en    = :name_en,
            name_es    = :name_es,
            decimals   = :decimals,
            fx_rate    = :fx_rate,
            is_enabled = :is_enabled,
            updated_at = NOW()
            WHERE id = :id
            RETURNING *`
        
        stmt, _ := db.PrepareNamed(query)
        if err := stmt.Get(&c_updated, c_body); err != nil {
            errRes := utils.ErrorResponse{}
            errRes.SetError("Failed to update item")
            errRes.Write(w, http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(c_updated)
    }
}

func DeleteCurrencyByID(db *sqlx.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")

        id, err := uuid.Parse(r.PathValue("id"))
        if err != nil {
            errRes := utils.ErrorResponse{}
            errRes.SetError("ID format is invalid")
            errRes.Write(w, http.StatusBadRequest)
            return
        }

        query := `DELETE FROM currencies WHERE id = $1`

        res, err := db.Exec(query, id)
        if err != nil {
            errRes := utils.ErrorResponse{}
            errRes.SetError("Failed to delete item")
            errRes.Write(w, http.StatusInternalServerError)
            return
        }

        // Check manually if any row was affected (SQL query always returns OK even if now rows were deleted)
        if rows, _ := res.RowsAffected(); rows == 0 {
            errRes := utils.ErrorResponse{}
            errRes.SetError("ID not found")
            errRes.Write(w, http.StatusNotFound)
            return
        }

        w.WriteHeader(http.StatusNoContent)
    }
}
