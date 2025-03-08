/*
Copyright © 2025 Francisco Brito

*/
package cmd

import (
	"os"
	"fmt"
	"github.com/spf13/cobra"
)



// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "javm",
    Short: "Gestor de versiones de Java",
    Long: `JAVM es una herramienta para administrar múltiples versiones de Java en tu sistema.

Permite instalar, cambiar y listar versiones del JDK fácilmente, 
simulando herramientas como 'nvm' para Node.js.`,
Run: func(cmd *cobra.Command, args []string) {
	fmt.Println("Bienvenido a JAVM, el gestor de versiones de Java.")
	fmt.Println("Usa 'javm --help' para ver los comandos disponibles.")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.javm.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


