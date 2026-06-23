package cmd

import (
	"fmt"
	"os"

	"github.com/Maxahup/cli-command-helper/internal/updater"
	"github.com/spf13/cobra"
)

const CurrentVersion = "v1.1.0"

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Actualiza el CLI a la última versión disponible desde GitHub",
	Long:  `Conecta con el repositorio público en GitHub para descargar e instalar de forma segura la última versión del ejecutable.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Buscando actualizaciones en GitHub...")

		err := updater.CheckAndApplyUpdate(CurrentVersion)
		if err != nil {
			fmt.Printf("Error durante la actualización: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
