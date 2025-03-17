# JAVM - Java Version Manager

![JAVM Logo](https://github.com/FranciscoJBrito/javm/assets/logo.png)

**JAVM** es un gestor de versiones de Java ligero y fácil de usar, inspirado en herramientas como `nvm` y `rbenv`. Permite instalar y cambiar entre múltiples versiones de Java en segundos, sin conflictos con versiones preinstaladas en el sistema.

---

## 🚀 Características
- 📌 Instalación y cambio de versiones de Java en segundos.
- 📥 Descarga automática desde **Adoptium**.
- 🔀 Alterna entre versiones con un solo comando.
- 🛠️ Configuración automática del entorno (`JAVA_HOME`, `PATH`).
- 🌍 Compatible con Linux, macOS y Windows.

---

## 📦 Instalación

JAVM se instala fácilmente ejecutando el siguiente script en la terminal:

```sh
curl -fsSL https://raw.githubusercontent.com/FranciscoJBrito/javm/main/install.sh | bash
```

### 🛠️ Instalación manual
Si prefieres instalarlo manualmente:

```sh
git clone https://github.com/FranciscoJBrito/javm.git
cd javm
chmod +x install.sh
./install.sh
```

⚠️ **Nota**: Asegúrate de tener permisos de escritura en `/usr/local/bin/`.

---

## 🚀 Uso Básico

### 📌 Ver las versiones de Java disponibles
```sh
javm list
```

### 📥 Instalar una versión específica
```sh
javm install 17
```

### 🔄 Cambiar de versión
```sh
javm use 17
```

### 🗑️ Eliminar una versión instalada
```sh
javm uninstall 8
```

### 🔎 Ver la versión actual en uso
```sh
javm current
```

---

## 🛠️ Contribuir
¡Las contribuciones son bienvenidas! Para colaborar:

1. **Forkea el repositorio**
2. **Crea una nueva rama** (`feature/nueva-funcionalidad`)
3. **Envía un Pull Request** 🚀

📌 Reporta bugs y sugiere mejoras en [GitHub Issues](https://github.com/FranciscoJBrito/javm/issues).

---

## 📜 Licencia
JAVM está bajo la licencia **MIT**. ¡Úsalo libremente! 🛠️

