package api

import (
    "net/http"

    "gocart/utils"
)

func answerInternalError(w http.ResponseWriter) {
    errRes := utils.ErrorResponse{}
    errRes.SetError("Unable to process the request. Please try again later. If the problem persists, please contact us.")
    errRes.Write(w, http.StatusInternalServerError)
}
