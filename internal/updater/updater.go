package updater

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/minio/selfupdate"
)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name        string `json:"name"`
		DownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func CheckAndApplyUpdate(currentVersion string) error {
	repoURL := "https://api.github.com/repos/Maxahup/cli-command-helper/releases/latest"

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(repoURL)
	if err != nil {
		return fmt.Errorf("no se pudo conectar con GitHub: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error al consultar GitHub (Status: %d)", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return fmt.Errorf("error al leer los datos de la versión: %v", err)
	}

	if release.TagName == currentVersion {
		fmt.Println("🎉 ¡Ya estás usando la última versión disponible (", currentVersion, ")!")
		return nil
	}

	fmt.Printf("🚀 Nueva versión detectada: %s (Actual: %s). Buscando binario compatible...\n", release.TagName, currentVersion)

	var downloadURL string
	expectedBinaryName := fmt.Sprintf("cli-helper-%s-%s", runtime.GOOS, runtime.GOARCH)

	for _, asset := range release.Assets {
		if asset.Name == expectedBinaryName {
			downloadURL = asset.DownloadURL
			break
		}
	}

	if downloadURL == "" {
		return fmt.Errorf("no se encontró un binario precompilado compatible para tu arquitectura (%s)", expectedBinaryName)
	}

	fmt.Println("Descargando actualización...")
	binaryResp, err := client.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("error al descargar el nuevo binario: %v", err)
	}
	defer binaryResp.Body.Close()

	err = selfupdate.Apply(binaryResp.Body, selfupdate.Options{})
	if err != nil {
		return fmt.Errorf("error crítico al aplicar la actualización en el disco: %v", err)
	}

	fmt.Printf("Actualizado correctamente - reinicia tu terminal para usar cli-helper %s\n", release.TagName)
	return nil
}
