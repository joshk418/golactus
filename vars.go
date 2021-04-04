package golactus

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Vars map[string]string

func GetVars(r *http.Request) Vars {
	return Vars(mux.Vars(r))
}

func (v Vars) TryInt(key string) (int, error) {
	return strconv.Atoi(v[key])
}

func (v Vars) Str(key string) string {
	return v[key]
}
