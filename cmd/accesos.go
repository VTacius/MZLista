package cmd

import (
	"fmt"
	"io"
	"os"
	"time"

	"xibalba.com/vtacius/MZLista/base"

	tm "github.com/buger/goterm"
	"github.com/spf13/cobra"
)

// accesosCmd represents the accesos command
var accesosCmd = &cobra.Command{
	Use:   "accesos",
	Short: "Muestra los inicios de sesión previos al envío de correo",
	Run: func(cmd *cobra.Command, args []string) {
		// Filtros para la información que se muestra
		ipaddress, _ := cmd.Flags().GetString("ip")
		username, _ := cmd.Flags().GetString("usuario")
		// Opciones más o menos normales de la operación
		fichero, _ := cmd.Flags().GetString("fichero")
		completo, _ := cmd.Flags().GetBool("completo")
		r, _ := cmd.Flags().GetInt("refresco")
		refresco := time.Duration(r)
		if completo {
			mostrarLista(fichero, ipaddress, username)
		} else {
			mostrarPager(fichero, refresco, ipaddress, username)
		}
	},
}

func init() {
	rootCmd.AddCommand(accesosCmd)
	accesosCmd.Flags().StringP("fichero", "f", "/var/log/zimbra.log", "Fichero a analizar")
	accesosCmd.Flags().BoolP("completo", "c", false, "Mostrar una lista fluida")
	accesosCmd.Flags().IntP("refresco", "r", 5, "Segundos de refresco")
	accesosCmd.Flags().String("ip", "", "IP de origen")
	accesosCmd.Flags().String("usuario", "", "Usuario")
}

func mostrarPager(fichero string, refresco time.Duration, ipaddress string, username string) {
	tm.Clear()
	for {
		tm.MoveCursor(1, 1)
		box := tm.NewBox(tm.Width(), tm.Height()-1, 0)

		listarAccesos(fichero, box, ipaddress, username)

		tm.Print(box.String())
		tm.Flush()
		time.Sleep(refresco * time.Second)
	}
}

func mostrarLista(fichero string, ipaddress string, username string) {
	salida := os.Stdout
	listarAccesos(fichero, salida, ipaddress, username)
}

func listarAccesos(fichero string, salida io.Writer, ipaddress string, username string) {
	parser := base.NewFichero(fichero)
	parser.Grepear()
	if ipaddress != "" {
		parser.FiltrarIP(ipaddress)
	}
	if username != "" {
		parser.FiltrarUsername(username)
	}
	hits := parser.Ordenar()
	for _, intentos := range hits {
		fmt.Fprint(salida, intentos.Mostrar(parser.Longitudes))
	}
}
