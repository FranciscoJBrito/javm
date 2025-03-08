package cmd

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/spf13/cobra"
    "github.com/FranciscoJBrito/jvm/internal/ui"
)

// useCmd representa el comando "use"
var useCmd = &cobra.Command{
    Use:   "use [version]",
    Short: "Cambia la versión de Java activa",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        version := args[0]
        jvmPath := filepath.Join(os.Getenv("HOME"), ".jvm", "versions", version)
        currentPath := filepath.Join(os.Getenv("HOME"), ".jvm", "current")

        if _, err := os.Stat(jvmPath); os.IsNotExist(err) {
            fmt.Println(ui.UseError.Render("❌ La versión especificada no está instalada."))
            return
        }

        // Si la versión ya está activa, evitar cambios innecesarios
        currentVersion, _ := os.Readlink(currentPath)
        if currentVersion == jvmPath {
            fmt.Println(ui.UseSuccess.Render("🔄 Java versión " + version + " ya está activa."))
            return
        }

        os.Remove(currentPath)
        os.Symlink(jvmPath, currentPath)

        fmt.Println(ui.UseSuccess.Render("✅ Ahora usando Java versión " + version))
    },
}

func init() {
    rootCmd.AddCommand(useCmd)
}
