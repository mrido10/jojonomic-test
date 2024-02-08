package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
)

func (a app) handleAPI(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		a.getListApi(w, r)
	case "POST":
		a.setListApi(w, r)
	}
}

type responseListAPI struct {
	ListApi []listApi `json:"list_api"`
}

type responseBody struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (a app) getListApi(w http.ResponseWriter, r *http.Request) {
	var list responseListAPI
	for _, v := range a.ListApi {
		list.ListApi = append(list.ListApi, v)
	}

	w.Header().Set("Content-Type", "application/json")
	result, err := json.Marshal(list)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(result)
}

type requestBodyApi listApi

func (a app) setListApi(w http.ResponseWriter, r *http.Request) {
	var body requestBodyApi
	w.Header().Set("Content-Type", "application/json")

	// read body request
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		setResult(w, http.StatusBadRequest, "Failed", "Request failed")
		return
	}

	// validate body request
	err = body.validateBodyRequest()
	if err != nil {
		setResult(w, http.StatusBadRequest, "Failed", err.Error())
		return
	}

	// set mock
	a.Route.HandleFunc("/api/"+body.Url, body.setMockApi).Methods(body.Method)

	// save to list
	body.Url = "localhost:8000/api/" + body.Url
	body.Response = ""
	a.ListApi[body.Url] = listApi(body)

	setResult(w, http.StatusOK, "Success save mock api", "")
}

func (rb requestBodyApi) validateBodyRequest() error {
	methodAllowed := map[string]bool{
		"GET":    true,
		"POST":   true,
		"PUT":    true,
		"DELETE": true,
	}
	if !methodAllowed[rb.Method] {
		return errors.New("method not allowed")
	}

	rgx, err := regexp.Compile(`[A-Za-z1-9]+$`)
	if err != nil {
		return err
	}

	if !rgx.MatchString(rb.Url) {
		return errors.New("url must alphanumeric")
	}

	response := make(map[string]interface{})
	err = json.Unmarshal([]byte(rb.Response), &response)
	if err != nil {
		return errors.New("response must be json")
	}

	return nil
}

func (rb requestBodyApi) setMockApi(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := make(map[string]interface{})
	err := json.Unmarshal([]byte(rb.Response), &response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(result)
}

func setResult(w http.ResponseWriter, code int, status, message string) {
	w.WriteHeader(code)
	result, _ := json.Marshal(responseBody{
		Status:  status,
		Message: message,
	})
	w.Write(result)
}
