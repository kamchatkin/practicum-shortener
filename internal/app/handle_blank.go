package app

import "net/http"

func BlankHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<script src= \n\"https://ajax.googleapis.com/ajax/libs/jquery/3.6.0/jquery.min.js\"> \n    </script> "))
}
