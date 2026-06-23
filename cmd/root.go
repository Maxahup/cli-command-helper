package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Maxahup/cli-command-helper/internal/ai"
	"github.com/Maxahup/cli-command-helper/internal/runner"
	"github.com/spf13/cobra"
)

var (
	flagOS      string
	flagVersion string
)

var rootCmd = &cobra.Command{
	Use:   "cli-helper [tu petición]",
	Short: "Asistente inteligente para la terminal",
	Long:  `Traduce lenguaje natural a comandos ejecutables de la shell inyectando contexto de OS.`,
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

		fmt.Printf("Analizando (OS: %s | Shell: %s)...\n", ctxIA.OS, ctxIA.Shell)

		aiClient := ai.NewClient("")

		comandoSugerido, err := aiClient.GetTerminalCommand(peticion, ctxIA)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		comandoLimpio := strings.TrimSpace(comandoSugerido)
		comandoLimpio = strings.ReplaceAll(comandoLimpio, "```sh", "")
		comandoLimpio = strings.ReplaceAll(comandoLimpio, "```bash", "")
		comandoLimpio = strings.ReplaceAll(comandoLimpio, "```", "")
		comandoLimpio = strings.TrimSpace(comandoLimpio)

		fmt.Println("\nComando sugerido:")
		fmt.Printf("   %s\n\n", comandoLimpio)

		fmt.Print("¿Deseas ejecutar este comando ahora mismo? [y/N]: ")

		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error al leer la respuesta: %v\n", err)
			os.Exit(1)
		}

		input = strings.TrimSpace(strings.ToLower(input))

		if input == "y" || input == "yes" {
			fmt.Println("Ejecutando comando...")
			fmt.Println("--------------------------------------------------")

			var cmdExec *exec.Cmd

			if ctxIA.OS == "windows" {
				if ctxIA.Shell == "powershell" {
					cmdExec = exec.Command("powershell", "-Command", comandoLimpio)
				} else {
					cmdExec = exec.Command("cmd", "/C", comandoLimpio)
				}
			} else {
				shell := ctxIA.Shell
				if shell == "unknown" || shell == "" {
					shell = "sh"
				}
				cmdExec = exec.Command(shell, "-c", comandoLimpio)
			}

			cmdExec.Stdout = os.Stdout
			cmdExec.Stderr = os.Stderr
			cmdExec.Stdin = os.Stdin

			if err := cmdExec.Run(); err != nil {
				fmt.Printf("\n❌ El comando finalizó con un error de ejecución: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("--------------------------------------------------")
			fmt.Println("Proceso terminado exitosamente.")
		} else {
			fmt.Println("Ejecución cancelada por el usuario.")
		}
	},
}

func init() {
	rootCmd.Flags().StringVarP(&flagOS, "os", "o", "", "Forzar un sistema operativo objetivo (windows, linux, darwin)")
	rootCmd.Flags().StringVarP(&flagVersion, "target-version", "v", "", "Especificar una versión de librería o entorno")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
