package ai

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

type AIClient struct {
	client *openai.Client
	model  string
}
type RequestContext struct {
	OS            string
	Shell         string
	TargetVersion string
}

func NewClient(modelName string) *AIClient {
	if modelName == "" {
		modelName = "qwen2.5-coder:7b"
	}

	config := openai.DefaultConfig("ollama")
	config.BaseURL = "http://localhost:11434/v1"

	return &AIClient{
		client: openai.NewClientWithConfig(config),
		model:  modelName,
	}
}

func (a *AIClient) GetTerminalCommand(prompt string, sysCtx RequestContext) (string, error) {
	ctx := context.Background()

	systemPrompt := fmt.Sprintf(
		"Eres un motor de traducción de lenguaje natural a comandos de terminal.\n\n"+
			"ENTORNO DE EJECUCIÓN ACTUAL:\n"+
			"- Sistema Operativo: %s\n"+
			"- Shell Activa: %s\n"+
			"- Versión/Herramienta Objetivo: %s\n\n"+
			"REGLAS ESTRICTAS DE RESPUESTA:\n"+
			"1. Devuelve ÚNICAMENTE la línea de comando ejecutable que resuelva la petición. No des múltiples opciones.\n"+
			"2. NO uses bloques de código de markdown (absolutamente prohibido usar ```bash ... ```).\n"+
			"3. NO des explicaciones, ni introducciones, ni comentarios post-comando, ni saludos.\n"+
			"4. Si el comando no se puede ejecutar de forma nativa en este OS/Shell, devuelve la alternativa exacta.",
		sysCtx.OS, sysCtx.Shell, sysCtx.TargetVersion,
	)

	resp, err := a.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: a.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.1,
		},
	)

	if err != nil {
		return "", fmt.Errorf("error al conectar con Ollama: %v", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no se recibió respuesta del modelo")
	}

	return resp.Choices[0].Message.Content, nil
}
