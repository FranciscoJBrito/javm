package cmd

import (
    "archive/tar"
    "archive/zip"
    "compress/gzip"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "runtime"
	"strings"

    "github.com/spf13/cobra"
    "github.com/FranciscoJBrito/jvm/internal/ui"
)

// installCmd representa el comando "install"
var installCmd = &cobra.Command{
    Use:   "install [version]",
    Short: "Instala una versión específica de Java",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        version := args[0]
        fmt.Println(ui.TitleStyle.Render("📥 Descargando Java versión " + version + "..."))

        osType := runtime.GOOS
		fmt.Println("Sistema operativo:", osType)

        // Verificar si la versión está disponible
        available, err := isVersionAvailable(version)
        if err != nil {
            fmt.Println(ui.ErrorStyle.Render("❌ Error al verificar la disponibilidad de la versión: " + err.Error()))
            return
        }
        if !available {
            fmt.Println(ui.ErrorStyle.Render("❌ La versión " + version + " no está disponible en Adoptium."))
            return
        }

        // Obtener URL de descarga
        url, err := buildDownloadURL(version)
        if err != nil {
            fmt.Println(ui.ErrorStyle.Render("❌ Error al construir la URL de descarga: " + err.Error()))
            return
        }

        // Descargar y extraer OpenJDK
        err = downloadAndInstallJDK(version, url, osType)
        if err != nil {
            fmt.Println(ui.ErrorStyle.Render("❌ Error al instalar JDK: " + err.Error()))
            return
        }

        fmt.Println(ui.SuccessStyle.Render("✅ Java versión " + version + " instalada correctamente."))
    },
}


func init() {
    rootCmd.AddCommand(installCmd)
}

// Verifica si una versión está disponible en Adoptium
func isVersionAvailable(version string) (bool, error) {
    url := "https://api.adoptium.net/v3/info/available_releases"
    resp, err := http.Get(url)
    if err != nil {
        return false, fmt.Errorf("fallo al consultar la API de Adoptium: %w", err)
    }
    defer resp.Body.Close()

    var data struct {
        AvailableReleases []int `json:"available_releases"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
        return false, fmt.Errorf("fallo al procesar la lista de versiones: %w", err)
    }

    for _, v := range data.AvailableReleases {
        if fmt.Sprintf("%d", v) == version {
            return true, nil
        }
    }
    return false, nil
}

// Construye la URL de descarga basada en la API de Adoptium
func buildDownloadURL(version string) (string, error) {
    osType := runtime.GOOS
    arch := runtime.GOARCH

    var osStr, archStr string

    switch osType {
    case "linux":
        osStr = "linux"
    case "darwin":
        osStr = "mac"
    case "windows":
        osStr = "windows"
    default:
        return "", fmt.Errorf("sistema operativo no soportado: %s", osType)
    }

    switch arch {
    case "amd64":
        archStr = "x64"
    case "arm64":
        archStr = "aarch64"
    default:
        return "", fmt.Errorf("arquitectura no soportada: %s", arch)
    }

    // Construir la URL de descarga
    url := fmt.Sprintf("https://api.adoptium.net/v3/binary/latest/%s/ga/%s/%s/jdk/hotspot/normal/eclipse",
        version, osStr, archStr)

    return url, nil
}

// Descargar y extraer OpenJDK
func downloadAndInstallJDK(version, url, osType string) error {
    installPath := filepath.Join(os.Getenv("HOME"), ".jvm", "versions", version)
	fmt.Println("Instalando en:", installPath)

    // Definir el nombre del archivo temporal según el SO
    tmpFile := filepath.Join(os.TempDir(), "jdk")
    if osType == "windows" {
        tmpFile += ".zip"
    } else {
        tmpFile += ".tar.gz"
    }
	fmt.Println("Archivo temporal:", tmpFile)

    // Descargar el archivo
    fmt.Println("📡 Descargando desde:", url)
    resp, err := http.Get(url)
    if err != nil {
        return fmt.Errorf("fallo al descargar JDK: %w", err)
    }
    defer resp.Body.Close()
	fmt.Println("Descarga completada.", resp)

    out, err := os.Create(tmpFile)
    if err != nil {
        return fmt.Errorf("fallo al crear archivo temporal: %w", err)
    }
    defer out.Close()

    _, err = io.Copy(out, resp.Body)
    if err != nil {
        return fmt.Errorf("fallo al guardar el archivo: %w", err)
    }

    fmt.Println("📦 Extrayendo JDK...")

    // Elegir el método correcto de extracción
    if osType == "windows" {
        err = extractZip(tmpFile, installPath)
    } else {
        err = extractTarGz(tmpFile, installPath)
    }

    if err != nil {
        return fmt.Errorf("fallo al extraer JDK: %w", err)
    }

    // Asegurar permisos de ejecución
    execPath := filepath.Join(installPath, "bin", "java")
    os.Chmod(execPath, 0755)

    os.Remove(tmpFile)
    fmt.Println(ui.SuccessStyle.Render("✅ JDK instalado en " + installPath))
    return nil
}

func extractTarGz(src, dest string) error {
    fmt.Println("📂 Iniciando extracción en:", dest)

    // Volver a abrir el archivo después de descargarlo
    file, err := os.Open(src)
    if err != nil {
        return fmt.Errorf("error al abrir el archivo comprimido: %w", err)
    }
    defer file.Close()

    gzipReader, err := gzip.NewReader(file)
    if err != nil {
        return fmt.Errorf("error al abrir gzip: %w", err)
    }
    defer gzipReader.Close()

    tarReader := tar.NewReader(gzipReader)

    // Crear el directorio de destino si no existe
    os.MkdirAll(dest, 0755)

    for {
        header, err := tarReader.Next()
        if err == io.EOF {
            break
        }
        if err != nil {
            return fmt.Errorf("error al leer el archivo tar: %w", err)
        }

        target := filepath.Join(dest, header.Name)

        // Prevenir ataques de path traversal
        if !strings.HasPrefix(target, filepath.Clean(dest)+string(os.PathSeparator)) {
            return fmt.Errorf("archivo fuera del directorio destino: %s", target)
        }

        switch header.Typeflag {
        case tar.TypeDir:
            fmt.Println("📁 Creando directorio:", target)
            if err := os.MkdirAll(target, 0755); err != nil {
                return fmt.Errorf("error al crear directorio: %w", err)
            }
        case tar.TypeReg:
            fmt.Println("📄 Extrayendo archivo:", target)
            if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
                return fmt.Errorf("error al crear directorios padre: %w", err)
            }

            outFile, err := os.Create(target)
            if err != nil {
                return fmt.Errorf("error al crear archivo de salida: %w", err)
            }
            defer outFile.Close()

            if _, err := io.Copy(outFile, tarReader); err != nil {
                return fmt.Errorf("error al copiar datos del archivo: %w", err)
            }
        default:
            fmt.Println("⚠️ Ignorando tipo de archivo no soportado:", header.Typeflag)
        }
    }

    fmt.Println("✅ Extracción completada correctamente en:", dest)
    return nil
}



// Extrae un archivo .zip a la ruta de instalación (solo Windows)
func extractZip(src, dest string) error {
    r, err := zip.OpenReader(src)
    if err != nil {
        return err
    }
    defer r.Close()

    os.MkdirAll(dest, 0755)

    for _, f := range r.File {
        fpath := filepath.Join(dest, f.Name)

        if f.FileInfo().IsDir() {
            os.MkdirAll(fpath, 0755)
            continue
        }

        rc, err := f.Open()
        if err != nil {
            return err
        }
        defer rc.Close()

        outFile, err := os.Create(fpath)
        if err != nil {
            return err
        }
        defer outFile.Close()

        _, err = io.Copy(outFile, rc)
        if err != nil {
            return err
        }
    }
    return nil
}