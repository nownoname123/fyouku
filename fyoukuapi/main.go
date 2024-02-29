package main

import (
	"fyoukuapi/router"
)

func main() {
	r := router.Router()

	err := r.Run("127.0.0.1:8005")
	if err != nil {
		return
	} // listen and serve on 0.0.0.0:8080
}
