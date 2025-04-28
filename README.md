# BlueNuke

**BlueNuke** es una herramienta de *hacking ético* desarrollada en Go para la detección, selección y prueba de seguridad sobre dispositivos Bluetooth cercanos.

> ⚡ **BlueNuke está desarrollado exclusivamente para uso ético y educativo.**  
> **Su uso en sistemas o dispositivos sin autorización explícita puede ser ilegal.**

---

## 🎯 Funcionalidades

- Escaneo en tiempo real de dispositivos Bluetooth.
- Guardado automático de dispositivos únicos en `dispositivos.txt`.
- Menú interactivo (navegación con flechas) para seleccionar dispositivos.
- Múltiples ataques disponibles:
  - **L2CAP Flood** *(solo en Linux)*: saturación de pings Bluetooth.
  - **Fake Pairing Request**: envío de solicitudes de emparejamiento falsas.

---

## ⚙️ Requisitos

- **Go** 1.20 o superior.
- **Adaptador Bluetooth** funcional.
- **Linux** (recomendado) o **MacOS**.
- Paquetes necesarios en Linux:
  - `bluez`
  - `bluetoothctl`
  - (Permisos de root para acceder a Bluetooth: `sudo`)

---

## 🚀 Instalación y uso rápido

1. Clona el repositorio:
```bash
git clone https://github.com/fervigz090/BlueNuke.git
cd BlueNuke
```

2. Instalar dependencias de Go
```bash
go mod tidy
```
3. Ejecutar BlueNuke
```bash
sudo go run bluenuke.go
```
4. Controles rápidos

[m] → Mostrar menú de dispositivos detectados.

[q] → Salir de BlueNuke.

🛡️ Uso Ético
BlueNuke ha sido creado únicamente para propósitos educativos y auditorías de seguridad Bluetooth en entornos controlados.
Nunca utilices esta herramienta en sistemas, dispositivos o redes ajenas sin autorización explícita.
El mal uso de BlueNuke puede ser considerado delito bajo las leyes de tu país.

👤 Autor
Desarrollado por Iván Fernández Rodríguez.

Contribuciones, mejoras o forks son bienvenidos, siempre que respeten el propósito ético del proyecto.

📜 Licencia
Este proyecto está licenciado bajo la Licencia MIT.