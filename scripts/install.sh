#!/bin/bash
echo "📥 Descargando JAVM CLI..."
curl -L -o javm https://github.com/FranciscoJBrito/javm/releases/download/v0.1.0/javm-linux
chmod +x javm
sudo mv javm /usr/local/bin/javm

if command -v javm &> /dev/null; then
    echo "✅ Instalación completa. Usa 'javm' en tu terminal."
else
    echo "❌ Error: La instalación falló. Verifica los permisos e intenta de nuevo."
fi
