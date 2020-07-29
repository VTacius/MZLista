package base

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/go-ldap/ldap/v3"
	"xibalba.com/vtacius/MZLista/utils"
)

// Acceso : Establece la conexión y operaciones LDAP
type Acceso struct {
	Cliente ldap.Client
	Base    string
	attrs   []string
	Err     error
	Datos   []Objeto
	Data    *ldap.SearchResult
}

// Buscar : Lista todos los objetos tales solicitados
func (acceso *Acceso) Buscar(filtro string, atributos []string) *Acceso {
	acceso.attrs = atributos
	peticion := ldap.NewSearchRequest(acceso.Base, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, filtro, atributos, nil)
	respuesta, err := acceso.Cliente.Search(peticion)
	if err != nil {
		acceso.Err = err
		return acceso
	}

	acceso.Data = respuesta
	return acceso
}

func obtenerObjeto(entrada *ldap.Entry, atributos []string) map[string]string {
	resultado := make(map[string]string)
	for _, clave := range atributos {
		valor := strings.TrimSpace(entrada.GetAttributeValue(clave))
		if valor != "" {
			resultado[clave] = valor
		}
	}
	return resultado
}

// Listar : Devuelve las entradas en un formato más accesible
func (acceso *Acceso) Listar() *Acceso {
	if acceso.Err != nil {
		return acceso
	}
	var resultado []Objeto
	for _, item := range acceso.Data.Entries {
		contenido := obtenerObjeto(item, acceso.attrs)
		resultado = append(resultado, Objeto{item.DN, contenido})
	}

	acceso.Datos = resultado
	return acceso
}

// FiltrarInactivos : Reduce la lista a sólo aquellos usuarios cuyo periodo de inactividad sea tal
func (acceso *Acceso) FiltrarInactivos(periodoInactividad float64) *Acceso {
	hoy := time.Now()
	var resultado []Objeto
	for item := range acceso.Datos {
		fecha := utils.Fechador(acceso.Datos[item].Atributos["zimbraLastLogonTimestamp"])
		if utils.RevisarIntervalo(periodoInactividad, hoy, fecha) == 1 {
			nuevo := acceso.Datos[item]
			nuevo.Atributos["zimbraLastLogonTimestamp"] = fecha.Format(time.RFC822)
			resultado = append(resultado, nuevo)
		}
	}

	acceso.Datos = resultado
	return acceso
}

func obtenerLongitudes(attrs *[]string, datos *[]Objeto) map[string]int {
	longitudes := make(map[string]int)
	for _, clave := range *attrs {
		longitudes[clave] = 0
	}

	for _, item := range *datos {
		for _, clave := range *attrs {
			l := len(item.Atributos[clave])
			if l > longitudes[clave] {
				longitudes[clave] = l
			}
		}
	}
	return longitudes
}

// ParaCSV : Produce una salida en CSV
func (acceso *Acceso) ParaCSV(salida io.Writer) {
	for _, item := range acceso.Datos {
		fmt.Fprintln(salida, item.Enumerar(acceso.attrs))
	}
}

// Imprimir : Muestra en pantalla el resultado
func (acceso *Acceso) Imprimir(salida io.Writer) {
	longitudes := obtenerLongitudes(&acceso.attrs, &acceso.Datos)
	for _, item := range acceso.Datos {
		fmt.Fprintln(salida, item.Tabular(acceso.attrs, longitudes))
	}
}
