package cmd

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/spf13/cobra"
    "github.com/FranciscoJBrito/jvm/internal/ui"
)

// listCmd representa el comando "list"
var listCmd = &cobra.Command{
    Use:   "list",
    Short: "Muestra las versiones de Java instaladas",
    Run: func(cmd *cobra.Command, args []string) {
        listInstalledVersions()
    },
}

func init() {
    rootCmd.AddCommand(listCmd)
}

// listInstalledVersions muestra todas las versiones instaladas de Java
func listInstalledVersions() {
    versionsDir := filepath.Join(os.Getenv("HOME"), ".jvm", "versions")
    currentLink := filepath.Join(os.Getenv("HOME"), ".jvm", "current")

    // Leer las versiones instaladas
    versions, err := os.ReadDir(versionsDir)
    if err != nil {
        fmt.Println(ui.ErrorStyle.Render("❌ Error al leer versiones instaladas: " + err.Error()))
        return
    }

    // Obtener la versión actualmente en uso
    currentVersion, err := os.Readlink(currentLink)
    if err != nil {
        currentVersion = "" // No hay una versión activa
    }

    fmt.Println(ui.TitleStyle.Render("📌 Versiones de Java instaladas:"))

    if len(versions) == 0 {
        fmt.Println(ui.ErrorStyle.Render("⚠️ No hay versiones instaladas. Usa `jvm install <versión>` para instalar una."))
        return
    }

    for _, v := range versions {
        versionName := v.Name()
        versionPath := filepath.Join(versionsDir, versionName)

        // Determinar si es la versión activa
        if currentVersion != "" && versionPath == currentVersion {
            fmt.Println(ui.SuccessStyle.Render("➡️ " + versionName + " (en uso)"))
        } else {
            fmt.Println("  " + versionName)
        }
    }
}
