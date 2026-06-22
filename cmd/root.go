package cmd

import (
	"fmt"
	"os"

	"github.com/Maxahup/cli-command-helper/internal/ai"
	"github.com/Maxahup/cli-command-helper/internal/runner"
	"github.com/spf13/cobra"
)

var (
	flagOS      string
	flagVersion string
)

var rootCmd = &cobra.Command{
	Use:   "ai-cli [tu petición]",
	Short: "Asistente inteligente para la terminal",
	Long:  `Traduce lenguaje natural a comandos ejecutables de la shell inyectando contexto de OS y librerías.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		peticion := args[0]

		contextoReal := runner.GetCurrentContext()

		ctxIA := ai.RequestContext{
			OS:            contextoReal.OS,
			Shell:         contextoReal.Shell,
			TargetVersion: "última estándar",
		}

		if flagOS != "" {
			ctxIA.OS = flagOS
		}
		if flagVersion != "" {
			ctxIA.TargetVersion = flagVersion
		}

		fmt.Printf("🤖 Analizando (OS: %s | Shell: %s)...\n", ctxIA.OS, ctxIA.Shell)

		aiClient := ai.NewClient("")

		comandoSugerido, err := aiClient.GetTerminalCommand(peticion, ctxIA)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Comando sugerido:")
		fmt.Printf("   %s\n\n", comandoSugerido)
	},
}

func init() {
	rootCmd.Flags().StringVarP(&flagOS, "os", "o", "", "Forzar un sistema operativo objetivo (windows, linux, darwin)")
	rootCmd.Flags().StringVarP(&flagVersion, "target-version", "v", "", "Especificar una versión de librería o entorno (ej: python3.9, docker-v2)")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
