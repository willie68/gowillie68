package measurement

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/samber/do"
)

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/", GetMetrics)
	return router
}

func GetMetrics(response http.ResponseWriter, request *http.Request) {
	ms := do.MustInvoke[*Service](nil)
	datas := ms.Datas()

	render.Status(request, http.StatusOK)
	render.JSON(response, request, datas)
}
