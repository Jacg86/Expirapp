package server

//Este archivo contiene la configuración del servidor Gin

import "time"

/*
# Config representa la configuración del servidor Gin
* Port: el puerto en el que se ejecuta el servidor
* ReadTimeout: el tiempo máximo para leer la solicitud completa, incluyendo el cuerpo
* WriteTimeout: el tiempo máximo para escribir la respuesta completa
* IdleTimeout: el tiempo máximo para esperar la próxima solicitud cuando keep-alives están habilitados
* Mode: el modo de ejecución de Gin (debug, release, test)
*/

type Config struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	Mode         string
}

//configuracion por defecto, de momento esa esta bien
func DefConfig() Config {
	return Config{
		Port:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Mode:         "debug",
	}
}
