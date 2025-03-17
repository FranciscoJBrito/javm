package cmd

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strings"

    "github.com/spf13/cobra"
    "github.com/FranciscoJBrito/javm/internal/ui"
)

// useCmd representa el comando "use"
var useCmd = &cobra.Command{
    Use:   "use [version]",
    Short: "Cambia la versión de Java en uso",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        version := args[0]
        switchJavaVersion(version)
    },
}

func init() {
    rootCmd.AddCommand(useCmd)
}

// switchJavaVersion cambia la versión de Java en uso
func switchJavaVersion(version string) {
    homeDir := os.Getenv("HOME")
    versionsDir := filepath.Join(homeDir, ".javm", "versions")
    currentPath := filepath.Join(homeDir, ".javm", "current")

    newVersionPath := filepath.Join(versionsDir, version)
    if _, err := os.Stat(newVersionPath); os.IsNotExist(err) {
        fmt.Println(ui.ErrorStyle.Render("❌ La versión especificada no está instalada."))
        return
    }

    // Actualizar el enlace simbólico ~/.javm/current
    os.Remove(currentPath)
    err := os.Symlink(newVersionPath, currentPath)
    if err != nil {
        fmt.Println(ui.ErrorStyle.Render("❌ Error al cambiar la versión de Java: " + err.Error()))
        return
    }

    // Dar permisos de ejecución a todos los binarios en la nueva versión
    binPath := filepath.Join(newVersionPath, "bin")
    os.Chmod(binPath, 0755)

    fmt.Println(ui.SuccessStyle.Render("✅ Ahora estás usando Java " + version))

    // Actualizar ~/.zshrc si es necesario y aplicar cambios en la sesión actual
    updateShellConfig(newVersionPath)
}


// updateShellConfig actualiza automáticamente JAVA_HOME y PATH en ~/.zshrc
func updateShellConfig(javaHome string) {
    homeDir := os.Getenv("HOME")
    shellConfig := filepath.Join(homeDir, ".zshrc") // Cambia a .bashrc si usas Bash

    configLine := `export JAVA_HOME=$(ls -d $HOME/.javm/current/*/ 2>/dev/null | head -n 1)`
    pathLine := `export PATH="$JAVA_HOME/bin:$PATH"`

    // Leer el contenido de ~/.zshrc
    content, err := os.ReadFile(shellConfig)
    if err != nil {
        fmt.Println(ui.ErrorStyle.Render("❌ Error al leer ~/.zshrc"))
        return
    }

    existingConfig := string(content)
    shouldAddConfig := !(strings.Contains(existingConfig, configLine) && strings.Contains(existingConfig, pathLine))

    if shouldAddConfig {
        // Agregar configuración si no existe
        fmt.Println(ui.WarningStyle.Render("⚠️ Configuración no encontrada en ~/.zshrc. Agregándola..."))

        newConfig := fmt.Sprintf("\n# Configuración automática de javm\n%s\n%s\n", configLine, pathLine)
        file, err := os.OpenFile(shellConfig, os.O_APPEND|os.O_WRONLY, 0644)
        if err != nil {
            fmt.Println(ui.ErrorStyle.Render("❌ Error al actualizar " + shellConfig))
            return
        }
        defer file.Close()

        file.WriteString(newConfig)
        fmt.Println(ui.SuccessStyle.Render("✅ Configuración de JAVA_HOME agregada en " + shellConfig))
    } else {
        fmt.Println(ui.SuccessStyle.Render("✅ Configuración ya existe en ~/.zshrc"))
    }

    // Ejecutar los export en la sesión actual para aplicar los cambios inmediatamente
    fmt.Println(ui.TitleStyle.Render("\n🔄 Aplicando cambios en la sesión actual..."))
    exec.Command("zsh", "-c", fmt.Sprintf("export JAVA_HOME='%s' && export PATH=\"$JAVA_HOME/bin:$PATH\"", javaHome)).Run()
}



