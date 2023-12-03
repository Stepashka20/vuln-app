package server

import "os"

func Init() {

	r := NewRouter()
	err := r.Run(":" + os.Getenv("SERVER_PORT"))
	if err != nil {
		panic(err)
	}
}
