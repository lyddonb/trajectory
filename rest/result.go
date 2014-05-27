package rest

import (
	"encoding/json"
	"net/http"
)

type Response map[string]interface{}

func (r *Response) toJson() ([]byte, error) {
	return json.Marshal(r)
}

func makeSuccessResponse(response interface{}) Response {
	return Response{
		"success": true,
		"result":  response,
	}
}

func makeErrorResponse(response interface{}, err string) Response {
	return Response{
		"success": false,
		"result":  response,
		"error":   err,
	}
}

func writeJsonResult(w http.ResponseWriter, result []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

func writeJsonError(w http.ResponseWriter, result []byte) {
	w.Header().Set("Content-Type", "application/json")
	http.Error(w, string(result), 500)
}

func SendJsonErrorResponse(w http.ResponseWriter, response interface{}, err string) {
	// TODO: Check for a debug flag.
	result := makeErrorResponse(response, err)
	jsonResult, _ := result.toJson()

	// TODO: Pass in a 500 or something.
	writeJsonError(w, jsonResult)
}

func SendJsonResponse(w http.ResponseWriter, response interface{}, err error) {
	if err != nil {
		SendJsonErrorResponse(w, response, err.Error())
		return
	}

	result := makeSuccessResponse(response)
	jsonResult, jsonErr := result.toJson()

	if jsonErr != nil {
		SendJsonErrorResponse(w, result, err.Error())
		return
	}

	writeJsonResult(w, jsonResult)
}
