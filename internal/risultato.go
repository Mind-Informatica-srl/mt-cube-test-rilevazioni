package mtcubetest

import (
	"fmt"
	"net/http"
)

type risultato struct {
	indice      int
	rilevazione Rilevazione
	response    http.Response
	esito       bool
	err         error
}

func (r risultato) Logga() string {
	esito := "OK"
	if !r.esito {
		esito = "ERR"
	}
	err := "rilevazione recapitata"
	if r.err != nil {
		err = r.err.Error()
	}
	return fmt.Sprintf("Rilevazione %d - %s - ID sensore: %s - Esito: %s", (r.indice + 1), esito, r.rilevazione.Position, err)
}
