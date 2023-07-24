package mtcubetest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

func NewCmd() (cmd *cobra.Command) {
	var position string
	var mtcubeEndpoint string
	var timeInterval string = "5"
	var callNumber string = "10"
	cmd = &cobra.Command{
		Use:   "mtcube [flags]",
		Short: "Start the MTCube Suite Program",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// recupero i parametri da castare
			var interval int64
			if interval, err = strconv.ParseInt(timeInterval, 10, 64); err != nil {
				return
			}
			interval = interval * 1000000000
			var callN int
			if callN, err = strconv.Atoi(callNumber); err != nil {
				return
			}
			// loggo l'inizio della simulazione
			fmt.Println("Avvio della simulazione di invio rilevazioni a MTCUBE")
			fmt.Printf("Verranno eseguite %d chiamate, ognuna ogni %s secondi a %s, usando %s come id del sensore\n", callN, timeInterval, mtcubeEndpoint, position)
			// istanzio gli indicatori del risultato
			var ok, nonok int
			// costruisco la url
			url := fmt.Sprintf("%s/api/v1/rilevazioni", mtcubeEndpoint)
			// per il numero di volte in callnumber
			for i := 0; i < callN; i++ {
				// istanzio il risultato
				risultato := risultato{
					indice: i,
				}
				// costruisco la rilevazione
				rilevazione := Rilevazione{
					Position:  position,
					ItemID:    strconv.Itoa(i),
					Timestamp: time.Now().Unix(),
				}
				risultato.rilevazione = rilevazione
				// invio la richiesta
				relSlilce := []Rilevazione{rilevazione}
				var relBytes []byte
				if relBytes, err = json.Marshal(&relSlilce); err != nil {
					return
				}
				var resp *http.Response
				if resp, err = http.Post(url, "application/json", bytes.NewBuffer(relBytes)); err != nil {
					risultato.err = err
					risultato.esito = false
				} else {
					risultato.response = *resp
					risultato.esito = true
				}
				// loggo l'esito
				fmt.Println(risultato.Logga())
				if risultato.esito {
					ok++
				} else {
					nonok++
				}
				// aspetto
				time.Sleep(time.Duration(interval))
			}
			// scrivo il risultato della simulazione
			fmt.Println("Fine della simulazione")
			fmt.Printf("%d rilevazioni inviate con successo. %d fallite.\n", ok, nonok)
			return nil
		},
		Args: cobra.ExactArgs(0),
	}
	cmd.Flags().StringVar(
		&mtcubeEndpoint,
		"mtcube-endpoint",
		os.Getenv("MTCUBE-ENDPOINT"), "L'indirizzo dell'installazione di MTCUBE a cui inviare le rilevazioni. Alternativamente valorizzare la variabile di ambiente MTCUBE-ENDPOINT. Parametro Obbligatorio. Esempio: https://mtcube.bcspeakers.com",
	)
	cmd.Flags().StringVar(
		&position,
		"position",
		os.Getenv("MTCUBE-POSITION"), "L'id del sensore da simulare nell'invio delle rilevazioni. Alternativamente valorizzare la variabile d'ambiente MTCUBE-POSITION. Parametro obbligatorio.",
	)
	cmd.Flags().StringVar(
		&timeInterval,
		"time-interval",
		os.Getenv("MTCUBE-TIME-INTERVAL"), "L'intervallo di tempo in secondi che il programma deve aspettare fra la simulazione di una rilevazione e quella successiva. Alternativamente valorizzare la variabile di ambiente MTCUBE-TIME-INTERVAL. Valore di default: 5",
	)
	cmd.Flags().StringVar(
		&callNumber,
		"call-number",
		os.Getenv("MTCUBE-CALL-NUMBER"), "Il numero di rilevazioni da inviare. Alternativamente valorizzare la variabile di ambiente MTCUBE-CALL-NUMBER. Valore di default: 10",
	)
	return

}
