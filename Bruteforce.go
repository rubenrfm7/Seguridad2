package main

import (
	"fmt"      
	"os/exec"  
	"sync"     
	"time"     
)

const alphabet = "abcdefghijklmnopqrstuvwxyz" 

func main() {
	cifrado := "archive.pdf.gpg" 
	maxLength := 4                // Longitud máxima de las combinaciones de contraseñas
	resultChan := make(chan string) // Canal para recibir resultados

	// Ajustar el número de trabajadores que generarán combinaciones
	numWorkers := 32
	var wg sync.WaitGroup // Grupo de espera para sincronizar goroutines

	startTime := time.Now() // Iniciar temporizador para medir el tiempo total de ejecución

	// Crear un canal para las combinaciones de contraseñas
	combinationChan := make(chan string)

	// Lanzar goroutines para probar combinaciones
	for i := 0; i < numWorkers; i++ {
		wg.Add(1) // Aumentar el contador del grupo de espera
		go func() {
			defer wg.Done() // Disminuir el contador al finalizar la goroutine
			for clave := range combinationChan { // Recibir combinaciones del canal
				// Mostrar la contraseña que se está probando
				fmt.Printf("Probando clave: %s\n", clave)

				// Intentar descifrar con la clave actual
				if result := tryDecrypt(clave, cifrado); result != "" {
					resultChan <- result // Enviar la clave encontrada al canal de resultados
					return 
				}
			}
		}()
	}

	// Lanzar goroutine para generar combinaciones de contraseñas
	go func() {
		defer close(combinationChan) // Cerrar el canal de combinaciones al finalizar
		generateCombinations(maxLength, combinationChan) // Generar combinaciones
	}()

	// Lanzar goroutine para esperar la finalización de los trabajadores
	go func() {
		wg.Wait() // Esperar a que todos los trabajadores terminen
		close(resultChan) // Cerrar el canal de resultados al finalizar
	}()

	// Esperar el resultado de la búsqueda de la clave
	for clave := range resultChan {
		elapsedTime := time.Since(startTime) // Calcular el tiempo transcurrido
		fmt.Printf("Clave encontrada: %s\n", clave) // Mostrar la clave encontrada
		fmt.Printf("Tiempo total transcurrido: %s\n", elapsedTime) // Mostrar el tiempo total
		return 
	}

	// Si no se encontró ninguna clave
	elapsedTime := time.Since(startTime) // Calcular el tiempo transcurrido
	fmt.Println("No se pudo encontrar la clave.") // Mensaje de error
	fmt.Printf("Tiempo total transcurrido: %s\n", elapsedTime) // Mostrar el tiempo total
}

// Función para generar combinaciones de contraseñas hasta una longitud máxima
func generateCombinations(maxLength int, combinationChan chan<- string) {
	for length := 1; length <= maxLength; length++ { // Iterar sobre longitudes de 1 a maxLength
		combinate("", length, combinationChan) // Generar combinaciones de la longitud actual
	}
}

// Función recursiva para generar combinaciones de contraseñas
func combinate(prefix string, length int, combinationChan chan<- string) {
	if length == 0 { // Caso base: si la longitud es 0, enviar el prefijo al canal
		combinationChan <- prefix
		return
	}

	// Iterar sobre cada carácter del alfabeto
	for _, char := range alphabet {
		// Llamar recursivamente añadiendo el carácter al prefijo
		combinate(prefix+string(char), length-1, combinationChan)
	}
}

// Función para intentar descifrar el archivo con una clave dada
func tryDecrypt(clave, cifrado string) string {
	cmd := exec.Command("gpg", "--batch", "--passphrase", clave, "-d", cifrado) // Comando para ejecutar gpg
	err := cmd.Run() // Ejecutar el comando
	if err == nil { // Si no hubo error, la clave es correcta
		return clave // Devolver la clave encontrada
	}
	return "" // Devolver vacío si la clave es incorrecta
}



