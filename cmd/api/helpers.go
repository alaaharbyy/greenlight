package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"maps"

	"github.com/julienschmidt/httprouter"
)

func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	// Use the json.MarshalIndent() function so that whitespace is added to the encoded
	// JSON. Here we use no line prefix ("") and tab indents ("\t") for each element.
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	/*
		for key, value := range headers {
			w.Header()[key] = value
		}
	*/
	maps.Copy(w.Header(), headers)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	err := json.NewDecoder(r.Body).Decode(dst)
	if err != nil {
		// If there is an error during decoding, start the triage...
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		switch {
		case errors.As(err, &syntaxError): //syntax err
			return fmt.Errorf("body contains badly-formated JSON (at character %d)", syntaxError.Offset)
			//$  curl -d '<?xml version="1.0" encoding="UTF-8"?><note><to>Alex</to></note>' localhost:4000/v1/movies
			//$ curl -d '{"title": "Moana", }' localhost:4000/v1/movies
		case errors.Is(err, io.ErrUnexpectedEOF): //syntax err Decode()
			return errors.New("body contains badly-formated JSON")

		case errors.As(err, &unmarshalTypeError): //JSON value is the wrong type for the target destination
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contans incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF): //Decode() will return this when the json object is empty
			return errors.New("body must not be empty")

			//$ curl -X POST localhost:4000/v1/movies

		//if we pass something that is not a non-nil pointer to Decode().
		case errors.As(err, &invalidUnmarshalError):
			panic(err) //panicking versus returning errors

		default:
			return err

		}

	}

	return nil
}
