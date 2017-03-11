package server


import "encoding/json"
import (
	"net/http"
	"reflect"
	"fmt"
	md "./models"
	"database/sql"
)

func RespondError(w http.ResponseWriter, code int, message string) {
	RespondJSON(w, code, map[string]string{"error": message})
}

func RespondJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Write(response)
}

func InspectEmpty(t interface{}) bool {
	s := reflect.ValueOf(t)
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		v := f.Interface()
		if  v == "" {
			return true
		}
	}
	return false
}

func CompareTypes(t1, t2 interface{})  {
	s1 := reflect.ValueOf(t1).Elem()
	//s1.Field(1).SetString("lalal")

	//fmt.Println("%s", s1.Field(1))
	s2 := reflect.ValueOf(t2).Elem()
	for i := 0; i < s1.NumField(); i++ {
		v1 := s1.Field(i)
		v2 := s2.Field(i)
		if  v1.String()=="" {
			v1.SetString(v2.String())
		}
	}
}

func PrintObject(t interface{})  {
	s := reflect.ValueOf(t)
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Printf("%d: %s %s = %v\n", i,
			typeOfT.Field(i).Name, f.Type(), f.Interface())
	}
}

func CheckDbErr(err error, w http.ResponseWriter) {
	switch err {
	case sql.ErrNoRows:
		RespondError(w, http.StatusNotFound, "Empty DB response")
	case md.UniqueError:
		RespondError(w, http.StatusConflict, err.Error())

	default:
		RespondError(w, http.StatusInternalServerError, err.Error())
	}
}