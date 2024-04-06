package ports

import "net/http"

type IHttp interface {
	//Создать поездку
	CreateTrip(w http.ResponseWriter, r *http.Request)
}
