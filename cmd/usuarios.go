/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"xibalba.com/vtacius/MZLista/base"
	"xibalba.com/vtacius/MZLista/utils"
)

var usuariosInactivosCmd = &cobra.Command{
	Use:   "inactivos",
	Short: "Lista de usuarios inactivos",
	Run: func(cmd *cobra.Command, args []string) {
		salida := os.Stdout
		filtro := "(&(ObjectClass=zimbraAccount)(zimbraLastLogonTimestamp=*))"

		paraCSV, _ := cmd.Flags().GetBool("csv")
		dominio, _ := cmd.Flags().GetString("dominio")
		meses, _ := cmd.Flags().GetInt("meses")
		periodoInactividad := float64(meses * 30 * 24 * 60 * 60)
		atributos := []string{"cn", "uid", "zimbraLastLogonTimestamp", "description"}

		baseDN := utils.ConstruirBase(dominio)
		url, usuario, contrasenia := utils.ParametrosAccesoLdap()

		conexion, err := base.Conectar(url, usuario, contrasenia)
		if err != nil {
			utils.Ruptura("Error al conectarse", err)
		}

		usuarios := base.Acceso{Base: baseDN, Cliente: conexion}
		usuarios.Buscar(filtro, atributos).Listar().FiltrarInactivos(periodoInactividad)
		if usuarios.Err != nil {
			utils.Ruptura("Error al listar usuarios", usuarios.Err)
		}

		// Imprime el resultado en pantalla en el formato requerido
		if paraCSV {
			usuarios.ParaCSV(salida)
		} else {
			usuarios.Imprimir(salida)
		}
	},
}

// usuariosCmd represents the usuarios command
var usuariosCmd = &cobra.Command{
	Use:   "usuarios",
	Short: "Lista los usuarios de correo",
	Run: func(cmd *cobra.Command, args []string) {
		salida := os.Stdout
		filtro := "(ObjectClass=zimbraAccount)"

		paraCSV, _ := cmd.Flags().GetBool("csv")
		dominio, _ := cmd.Flags().GetString("dominio")
		atributos, _ := cmd.Flags().GetStringArray("atributos")

		baseDN := utils.ConstruirBase(dominio)
		url, usuario, contrasenia := utils.ParametrosAccesoLdap()

		conexion, err := base.Conectar(url, usuario, contrasenia)
		if err != nil {
			utils.Ruptura("Error al conectarse", err)
		}

		usuarios := base.Acceso{Base: baseDN, Cliente: conexion}
		usuarios.Buscar(filtro, atributos).Listar()
		if usuarios.Err != nil {
			utils.Ruptura("Error al listar usuarios", usuarios.Err)
		}

		// Imprime el resultado en pantalla en el formato requerido
		if paraCSV {
			usuarios.ParaCSV(salida)
		} else {
			usuarios.Imprimir(salida)
		}

	},
}

func init() {
	rootCmd.AddCommand(usuariosCmd)
	// Todos los comandos necesitan dominio y formato de salida
	usuariosCmd.PersistentFlags().Bool("csv", false, "Muestra el resultado como CSV")
	usuariosCmd.PersistentFlags().StringP("dominio", "d", "sv", "Dominio sobre el cual buscar")

	usuariosCmd.Flags().StringArrayP("atributos", "a", []string{"uid", "displayName"}, "Atributos a buscar")

	// Este es un subcomando del subcomando
	usuariosCmd.AddCommand(usuariosInactivosCmd)
	usuariosInactivosCmd.Flags().Int("meses", 6, "Dias de inactividad")
}
