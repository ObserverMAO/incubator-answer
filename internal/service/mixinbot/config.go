package mixinbot

type MixinBotConfig struct {
	ClientID          string `json:"client_id" mapstructure:"client_id" yaml:"client_id"`
	SessionID         string `json:"session_id" mapstructure:"session_id" yaml:"session_id"`
	ServerPublicKey   string `json:"server_public_key" mapstructure:"server_public_key" yaml:"server_public_key"`
	SessionPrivateKey string `json:"session_private_key" mapstructure:"session_private_key" yaml:"session_private_key"`
	SpendKey          string `json:"spend_key" mapstructure:"spend_key" yaml:"spend_key"`
}
