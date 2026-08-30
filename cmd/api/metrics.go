package main

import (
	"fmt"
	"net/http"
)

func (app *application) getMetricsVolumeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Metrics by Volume Handler")
}
