package routes

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

const jsonContentType = "application/json"

func writeJSON(writer http.ResponseWriter, v interface{}) {
	body, err := marshallJSON(v)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = writer.Write(body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func marshallJSON(v interface{}) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func unmarkshallJSONBody(body io.ReadCloser, v interface{}) (err error) {
	return _unmarkshallJSONBody(body, v, func(b []byte) {})
}

func _unmarkshallJSONBody(body io.ReadCloser, v interface{}, callback func(b []byte)) (err error) {
	var b []byte
	b, err = ioutil.ReadAll(body)
	callback(b)
	defer func() {
		cerr := body.Close()
		if err == nil {
			err = cerr
		}
	}()
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, v)
	return err
}
