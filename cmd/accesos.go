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
		fichero, _ := cmd.Flags().GetString("fichero")
		completo, _ := cmd.Flags().GetBool("completo")
		r, _ := cmd.Flags().GetInt("refresco")
		refresco := time.Duration(r)
		if completo {
			mostrarLista(fichero)
		} else {
			mostrarPager(fichero, refresco)
		}
	},
}

func init() {
	rootCmd.AddCommand(accesosCmd)
	accesosCmd.Flags().StringP("fichero", "f", "/var/log/zimbra.log", "Fichero a analizar")
	accesosCmd.Flags().BoolP("completo", "c", false, "Mostrar una lista fluida")
	accesosCmd.Flags().IntP("refresco", "r", 5, "Segundos de refresco")
}

func mostrarPager(fichero string, refresco time.Duration) {
	tm.Clear()
	for {
		tm.MoveCursor(1, 1)
		box := tm.NewBox(tm.Width(), tm.Height()-1, 0)

		listarAccesos(fichero, box)

		tm.Print(box.String())
		tm.Flush()
		time.Sleep(refresco * time.Second)
	}
}

func mostrarLista(fichero string) {
	salida := os.Stdout
	listarAccesos(fichero, salida)
}

func listarAccesos(fichero string, salida io.Writer) {
	parser := base.NewFichero(fichero)
	hits := parser.Grepear().Ordenar()
	for _, intentos := range hits {
		fmt.Fprint(salida, intentos.Mostrar(parser.Longitudes))
	}
}
