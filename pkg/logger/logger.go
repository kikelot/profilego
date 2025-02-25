package logger

import (
	"log"
	"os"
)

// Logger personalizado
var Logger *log.Logger

// InitLogger inicializa el logger
func InitLogger() {
	// Crear archivo de logs
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error al abrir el archivo de logs: %v", err)
	}

	// Inicializar logger
	Logger = log.New(file, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile)
	log.Println("✅ Logger inicializado")
}

// LogInfo imprime logs de nivel informativo
func LogInfo(msg string) {
	Logger.Println("INFO: " + msg)
}

// LogError imprime logs de errores
func LogError(msg string) {
	Logger.Println("ERROR: " + msg)
}
