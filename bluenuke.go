package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"tinygo.org/x/bluetooth"
)

var (
	adapter            = bluetooth.DefaultAdapter
	dispositivosVistos = make(map[string]bool) // Mapa para no repetir dispositivos
)

func main() {
	printBanner()

	// Inicializar adaptador Bluetooth
	must("habilitar el adaptador Bluetooth", adapter.Enable())

	fmt.Println("[*] Escaneando dispositivos Bluetooth... (Ctrl+C para salir)")

	// Manejar interrupciones (Ctrl+C)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("\n[!] Interrupción detectada. Deteniendo escaneo...")
		adapter.StopScan()
		os.Exit(0)
	}()

	// Empezar escaneo
	err := adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
		address := device.Address.String()
		name := device.LocalName()
		timestamp := time.Now().Format("15:04:05") // Formato de hora HH:MM:SS

		// Si no lo hemos visto antes, lo agregamos
		if !dispositivosVistos[address] {
			dispositivosVistos[address] = true

			fmt.Printf("[%s] [+] Nuevo dispositivo: %s - %s (RSSI: %d)\n", timestamp, address, name, device.RSSI)

			// Guardar en archivo
			saveDevice(timestamp, address, name, device.RSSI)
		}
	})
	must("escanear", err)

	// Mantener corriendo
	for {
		time.Sleep(1 * time.Second)
	}
}

func saveDevice(timestamp, address, name string, rssi int16) {
	file, err := os.OpenFile("dispositivos.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("[-] Error al abrir el archivo: %s\n", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	line := fmt.Sprintf("[%s] %s - %s (RSSI: %d)\n", timestamp, address, name, rssi)
	writer.WriteString(line)
	writer.Flush()
}

func must(action string, err error) {
	if err != nil {
		panic(fmt.Sprintf("[-] Error al %s: %s", action, err))
	}
}

func printBanner() {
	banner := `
                                                                                                                
      _____    ____         ____   ____      ______  _____   ______    ____   ____  ____    ____       ______   
 ___|\     \  |    |       |    | |    | ___|\     \|\    \ |\     \  |    | |    ||    |  |    |  ___|\     \  
|    |\     \ |    |       |    | |    ||     \     \\\    \| \     \ |    | |    ||    |  |    | |     \     \ 
|    | |     ||    |       |    | |    ||     ,_____/|\|    \  \     ||    | |    ||    | /    // |     ,_____/|
|    | /_ _ / |    |  ____ |    | |    ||     \--'\_|/ |     \  |    ||    | |    ||    |/ _ _//  |     \--'\_|/
|    |\    \  |    | |    ||    | |    ||     /___/|   |      \ |    ||    | |    ||    |\    \'  |     /___/|  
|    | |    | |    | |    ||    | |    ||     \____|\  |    |\ \|    ||    | |    ||    | \    \  |     \____|\ 
|____|/____/| |____|/____/||\___\_|____||____ '     /| |____||\_____/||\___\_|____||____|  \____\ |____ '     /|
|    /     || |    |     ||| |    |    ||    /_____/ | |    |/ \|   ||| |    |    ||    |   |    ||    /_____/ |
|____|_____|/ |____|_____|/ \|____|____||____|     | / |____|   |___|/ \|____|____||____|   |____||____|     | /
  \(    )/      \(    )/       \(   )/    \( |_____|/    \(       )/      \(   )/    \(       )/    \( |_____|/ 
   '    '        '    '         '   '      '    )/        '       '        '   '      '       '      '    )/    
                                                '                                                         '     
                  Developed by Ivan Fernandez Rodriguez
`
	fmt.Println(banner)
}
