package httpjson

import (
	"encoding/json"
	"net/http"

	"github.com/CustodiaJS/custodiajs-core/static"
	"github.com/CustodiaJS/custodiajs-core/types"
)

func responseWrite(contentType types.HttpRequestContentType, w http.ResponseWriter, rpcd *ResponseCapsle) (int, error) {
	// Der Passende Typ wird festgelegt
	switch contentType {
	case static.HTTP_CONTENT_CBOR:
		w.Header().Set("Content-Type", "application/cbor")
	default:
		w.Header().Set("Content-Type", "application/json")
	}

	// Die Antwort wird gebaut
	var response *RPCResponse
	if rpcd.Data != nil {
		w.WriteHeader(http.StatusOK)
		response = &RPCResponse{Result: "ok", Error: nil, Data: rpcd.Data}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		response = &RPCResponse{Result: "failed", Error: &rpcd.Error, Data: nil}
	}

	// Das Antwortpaket wird ferigestellt
	var bytedResponse []byte
	var err error
	switch contentType {
	case static.HTTP_CONTENT_CBOR:
		bytedResponse, err = json.Marshal(response)
	default:
		bytedResponse, err = json.Marshal(response)
	}

	// Es wird geprüft ob ein Fehler aufgetreten ist
	if err != nil {
		panic(err)
	}

	// Das Paket wird gesendet
	return w.Write(bytedResponse)
}

func errorResponse(contentType types.HttpRequestContentType, w http.ResponseWriter, s string) (int, error) {
	return responseWrite(contentType, w, &ResponseCapsle{Error: s})
}
