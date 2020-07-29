package base

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

//IntentoAcceso : Información extraida de los registros de acceso
type IntentoAcceso struct {
	usuario   string
	hostname  string
	ipaddress string
}

// Hits : Lleva el control de acceso
type Hits struct {
	usuario   string
	hostname  string
	ipaddress string
	intento   int
}

// Mostrar : Formato para representar Hits como cadena
func (hits *Hits) Mostrar(longitudes Longitudes) string {
	var cadena strings.Builder
	cadena.WriteString(fmt.Sprintf("%*d ", 8, hits.intento))
	cadena.WriteString(fmt.Sprintf("%-*s ", longitudes.usuario, hits.usuario))
	cadena.WriteString(fmt.Sprintf("%-*s ", longitudes.hostname, hits.hostname))
	cadena.WriteString(fmt.Sprintf("%15s ", hits.ipaddress))
	cadena.WriteString("\n")
	return cadena.String()
}

// Longitudes : Guarda las longitudes a usar
type Longitudes struct {
	usuario  int
	hostname int
}

// Fichero : Encapsula las operaciones de ficheros
type Fichero struct {
	fichero    *os.File
	intentos   map[IntentoAcceso]int
	Longitudes Longitudes
}

// NewFichero : Inicializa un fichero
func NewFichero(ruta string) *Fichero {
	fichero, err := os.Open(ruta)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s %v\n", "Error abriendo fichero", err)
	}
	return &Fichero{fichero: fichero}
}

// Grepear : Busca las lineas coíncidentes y las pone en un formato adecuado
func (parser *Fichero) Grepear() *Fichero {
	defer parser.fichero.Close()
	intentos := make(map[IntentoAcceso]int)

	scanner := bufio.NewScanner(parser.fichero)
	for scanner.Scan() {
		texto := scanner.Text()
		// Grepeamos
		if indice := strings.Index(texto, "sasl_username"); indice > 0 {
			// Procesamos la linea y la ponemos en el map
			carga := SeleccionarCargaUtil(texto)

			usuario := ExtraerUsuario(carga[2])
			hostname, ipaddress := ParsearDireccion(carga[0])

			// Guardamos lo requerido
			intentos[IntentoAcceso{usuario, hostname, ipaddress}]++
		}
	}
	parser.intentos = intentos
	return parser
}

// FiltrarIP : Pues que filtra los intentos de conexión por IP
func (parser *Fichero) FiltrarIP(ipaddress string) *Fichero {
	intentos := make(map[IntentoAcceso]int)
	for cliente, intento := range parser.intentos {
		if cliente.ipaddress == ipaddress {
			intentos[cliente] = intento
		}
	}

	parser.intentos = intentos
	return parser
}

// FiltrarUsername : Pues que filtra los intentos de conexión por IP
func (parser *Fichero) FiltrarUsername(username string) *Fichero {
	intentos := make(map[IntentoAcceso]int)
	for cliente, intento := range parser.intentos {
		if cliente.usuario == username {
			intentos[cliente] = intento
		}
	}

	parser.intentos = intentos
	return parser
}

// Ordenar : Cambia el mapa a una lista ordenada
func (parser *Fichero) Ordenar() (hits []Hits) {
	var longitudes Longitudes

	// Habrá que cambiar el map a un slice de struct. Aprovechamos para encontrar las longitudes
	for acceso, intentos := range parser.intentos {
		item := Hits{acceso.usuario, acceso.hostname, acceso.ipaddress, intentos}
		if len(item.usuario) > longitudes.usuario {
			longitudes.usuario = len(item.usuario)
		}
		if len(item.hostname) > longitudes.hostname {
			longitudes.hostname = len(item.hostname)
		}
		hits = append(hits, item)
	}

	// Ahora si, podemos ordenar
	sort.Slice(hits, func(i, j int) bool {
		return hits[i].intento > hits[j].intento
	})

	parser.Longitudes = longitudes
	return
}

// SeleccionarCargaUtil : Toma la parte del registro que contiene la información a parsear
func SeleccionarCargaUtil(raw string) []string {
	indice := strings.Index(raw, "client=")
	indice += 7
	return strings.Fields(raw[indice:])
}

// ExtraerUsuario : Extrae el usuario que realiza la conexión
func ExtraerUsuario(raw string) (usuario string) {
	indice := strings.Index(raw, "=") + 1
	usuario = raw[indice:]

	return
}

// ParsearDireccion : Toma el usuario que realizó el acceso
func ParsearDireccion(raw string) (direccion, ipaddress string) {
	inicio := strings.Index(raw, "[")
	direccion = raw[:inicio]
	ipaddress = strings.Trim(raw[inicio+1:], "],")

	return
}
