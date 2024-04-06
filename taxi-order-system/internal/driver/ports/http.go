package ports

import "net/http"

type IHttp interface {
	//Захватить заявку
	CaptureTrip(w http.ResponseWriter, r *http.Request)
	//На позиции
	OnPosition(w http.ResponseWriter, r *http.Request)
	//Старт поездки
	Started(w http.ResponseWriter, r *http.Request)
	//Конец поездки
	Ended(w http.ResponseWriter, r *http.Request)
}
