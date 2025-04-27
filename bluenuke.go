package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/manifoldco/promptui"
	"tinygo.org/x/bluetooth"
)

var (
	adapter              = bluetooth.DefaultAdapter
	dispositivosVistos   = make(map[string]bool) // Para evitar repetidos
	listaDispositivos    = []string{}            // Para mostrar en menú
	addressToDeviceInfo  = make(map[string]DeviceInfo)
	imprimirDispositivos = true // Controlar impresión en pantalla
)

type DeviceInfo struct {
	Name string
	RSSI int16
	Time string
}

func main() {
	printBanner()

	must("habilitar el adaptador Bluetooth", adapter.Enable())

	fmt.Println("[*] Escaneando dispositivos Bluetooth...")
	fmt.Println("[m] Mostrar menú de dispositivos | [q] Salir")

	// Manejar interrupciones Ctrl+C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\n[!] Interrupción detectada. Deteniendo BlueNuke...")
		adapter.StopScan()
		os.Exit(0)
	}()

	// Escuchar entradas de teclado
	go listenForKeys()

	// Empezar escaneo
	err := adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
		address := device.Address.String()
		name := device.LocalName()
		timestamp := time.Now().Format("15:04:05")

		if !dispositivosVistos[address] {
			dispositivosVistos[address] = true

			info := fmt.Sprintf("%s - %s (RSSI: %d)", address, name, device.RSSI)
			listaDispositivos = append(listaDispositivos, info)

			addressToDeviceInfo[address] = DeviceInfo{
				Name: name,
				RSSI: device.RSSI,
				Time: timestamp,
			}

			if imprimirDispositivos {
				fmt.Printf("[%s] [+] Nuevo dispositivo: %s - %s (RSSI: %d)\n", timestamp, address, name, device.RSSI)
			}

			saveDevice(timestamp, address, name, device.RSSI)
		}
	})
	must("escanear", err)

	// Mantener programa vivo
	for {
		time.Sleep(1 * time.Second)
	}
}

func listenForKeys() {
	reader := bufio.NewReader(os.Stdin)
	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "m":
			showDeviceMenu()
		case "q":
			fmt.Println("[*] Saliendo de BlueNuke...")
			adapter.StopScan()
			os.Exit(0)
		default:
			fmt.Println("[!] Comando desconocido. Usa [m] para menú o [q] para salir.")
		}
	}
}

func showDeviceMenu() {
	imprimirDispositivos = false                   // Pausar impresión
	defer func() { imprimirDispositivos = true }() // Reanudar impresión al salir

	if len(listaDispositivos) == 0 {
		fmt.Println("[!] No se han detectado dispositivos todavía.")
		return
	}

	prompt := promptui.Select{
		Label: "Selecciona un dispositivo",
		Items: listaDispositivos,
		Size:  10,
	}

	index, _, err := prompt.Run()
	if err != nil {
		fmt.Printf("[-] Error al seleccionar: %v\n", err)
		return
	}

	selected := listaDispositivos[index]
	fmt.Println("\n[*] Dispositivo seleccionado:")
	fmt.Println(selected)

	fmt.Println("[!] (Próximamente: funcionalidades de ataque ético...)")
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
