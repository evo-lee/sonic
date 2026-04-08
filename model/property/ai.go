package property

import "reflect"

var (
	AIProvider = Property{
		DefaultValue: "anthropic",
		KeyValue:     "ai_provider",
		Kind:         reflect.String,
	}
	AIAPIKey = Property{
		DefaultValue: "",
		KeyValue:     "ai_api_key",
		Kind:         reflect.String,
	}
	AIModel = Property{
		DefaultValue: "",
		KeyValue:     "ai_model",
		Kind:         reflect.String,
	}
	// AIBaseURL is used by openai-compatible providers (OpenAI, Ollama, etc.)
	AIBaseURL = Property{
		DefaultValue: "",
		KeyValue:     "ai_base_url",
		Kind:         reflect.String,
	}
)
