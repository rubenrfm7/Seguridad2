package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"
)

func main() {
	cifrado := "archive.pdf.gpg" // Archivo cifrado a descifrar
	passwordFile := "hola.txt" // Archivo que contiene las contraseñas
	resultChan := make(chan string) // Canal para recibir resultados

	numWorkers := 32
	var wg sync.WaitGroup // Grupo de espera para sincronización de goroutines
	passwordSet := make(map[string]struct{}) // Conjunto para almacenar contraseñas únicas
	var mu sync.Mutex // Mutex para proteger el acceso al conjunto

	startTime := time.Now() // Iniciar temporizador para medir el tiempo total de ejecución

	// Cargar contraseñas en el conjunto
	for clave := range readPasswords(passwordFile) {
		mu.Lock()
		passwordSet[clave] = struct{}{} // Agregar la contraseña al conjunto
		mu.Unlock()
	}

	// Lanzar goroutines para probar contraseñas
	for i := 0; i < numWorkers; i++ {
		wg.Add(1) // Aumentar el contador del grupo de espera
		go func() {
			defer wg.Done() // Disminuir el contador al finalizar
			mu.Lock()
			for clave := range passwordSet { // Recibir contraseñas del conjunto
				delete(passwordSet, clave) // Eliminar la contraseña del conjunto para evitar repeticiones
				mu.Unlock() // Liberar el mutex antes de intentar descifrar

				// Mostrar la contraseña que se está probando
				fmt.Printf("Probando clave: %s\n", clave)

				// Intentar descifrar con la clave actual
				if result := tryDecrypt(clave, cifrado); result != "" {
					resultChan <- result // Enviar la clave encontrada al canal de resultados
					return // Terminar el goroutine si se encuentra la clave
				}

				mu.Lock() // Bloquear nuevamente para la siguiente iteración
			}
			mu.Unlock()
		}()
	}

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
		return // Salir de la función principal
	}

	// Si no se encontró ninguna clave
	elapsedTime := time.Since(startTime) // Calcular el tiempo transcurrido
	fmt.Println("No se pudo encontrar la clave.") // Mensaje de error
	fmt.Printf("Tiempo total transcurrido: %s\n", elapsedTime) // Mostrar el tiempo total
}

// Función para leer contraseñas desde un archivo
func readPasswords(filePath string) <-chan string {
	ch := make(chan string) // Canal para enviar contraseñas

	go func() {
		defer close(ch) // Cerrar el canal al finalizar
		file, err := os.Open(filePath) // Abrir el archivo
		if err != nil {
			fmt.Println("Error al abrir el archivo:", err)
			return
		}
		defer file.Close() // Asegurarse de cerrar el archivo al final

		scanner := bufio.NewScanner(file) // Crear un scanner para leer el archivo línea por línea
		for scanner.Scan() { // Leer cada línea
			ch <- scanner.Text() // Enviar la línea (contraseña) al canal
		}

		if err := scanner.Err(); err != nil { // Comprobar si hubo errores al leer
			fmt.Println("Error al leer el archivo:", err)
		}
	}()

	return ch // Devolver el canal
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

