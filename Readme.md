# cli-command-helper (`cli-helper`)

Un asistente inteligente de línea de comandos (CLI) de alto rendimiento escrito en **Go**. Traduce instrucciones en lenguaje natural a comandos puros de la terminal de manera **100% local, privada y gratuita** integrándose con **Ollama**.

El sistema analiza automáticamente tu entorno para garantizar que la sintaxis devuelta sea totalmente compatible con tu sistema operativo y shell actuales.

---

## Características Principales

* **Inyección de Contexto Dinámica:** Detecta automáticamente el Sistema Operativo (`darwin`, `linux`, `windows`) y la Shell activa (`zsh`, `bash`, `powershell`), adaptando las respuestas de la IA al instante.
* **Ejecución Interactiva Segura:** El CLI nunca ejecuta código a ciegas. Te presenta el comando sugerido y solicita una confirmación explícita `[y/N]` antes de lanzarlo al sistema.
* **Soporte Universal Estático:** Compilado con `CGO_ENABLED=0`, generando binarios independientes que no dependen de librerías externas del sistema.
* **Auto-Actualizaciones Integradas (`update`):** Descarga e instala en caliente las nuevas versiones directamente desde la sección de Releases de GitHub.

---

## Prerrequisitos Básicos

Para que la herramienta funcione localmente en tu máquina, necesitas tener configurado el motor de IA:

1. Descarga e instala **[Ollama](https://ollama.com/)**.
2. Descarga el modelo especializado en desarrollo de software corriendo en tu terminal:
   ```bash
   ollama pull qwen2.5-coder:7b