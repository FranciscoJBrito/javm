#!/bin/bash
echo "📥 Descargando JVM CLI..."
curl -L -o jvm https://github.com/FranciscoJBrito/jvm/releases/download/v0.1.0/jvm-linux
chmod +x jvm
sudo mv jvm /usr/local/bin/jvm

if command -v jvm &> /dev/null; then
    echo "✅ Instalación completa. Usa 'jvm' en tu terminal."
else
    echo "❌ Error: La instalación falló. Verifica los permisos e intenta de nuevo."
fi
