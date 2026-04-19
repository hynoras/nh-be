package secret

import (
	"context"
	"fmt"
	"log"
	"os"

	vault "github.com/hashicorp/vault/api"
)

func NewVaultClient() *vault.Client {
	config := vault.DefaultConfig()
	config.Address = os.Getenv("VAULT_ADDR")

	client, err := vault.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func AuthenticateVault(client *vault.Client) {
	data := map[string]interface{}{
		"role_id":   os.Getenv("VAULT_ROLE_ID"),
		"secret_id": os.Getenv("VAULT_SECRET_ID"),
	}

	secret, err := client.Logical().Write("auth/approle/login", data)
	if err != nil {
		log.Fatal(err)
	}

	client.SetToken(secret.Auth.ClientToken)
}

// GetSecret retrieves all key-value pairs from a KV v2 secret path.
// The path should NOT include the "secret/data/" prefix — the KV v2 client handles that.
// Example: GetSecret(client, "noheir/dev/db")
func GetSecret(client *vault.Client, path string) (map[string]interface{}, error) {
	kv := client.KVv2("secret")

	secret, err := kv.Get(context.Background(), path)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret at %s: %w", path, err)
	}

	return secret.Data, nil
}

// MustGetSecretValue extracts a string value from a Vault secret data map.
// Panics if the key is missing or not a string, matching the fail-fast behavior of env.MustEnv.
func MustGetSecretValue(data map[string]interface{}, key string) string {
	val, ok := data[key]
	if !ok {
		panic(fmt.Sprintf("missing required vault secret key: %s", key))
	}

	str, ok := val.(string)
	if !ok {
		panic(fmt.Sprintf("vault secret key %s is not a string", key))
	}

	return str
}