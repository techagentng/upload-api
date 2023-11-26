package config

import (
	"context"
	"fmt"
	"log"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Debug                        bool   `envconfig:"debug"`
	Port                         int    `envconfig:"port"`
	PostgresHost                 string `envconfig:"postgres_host"`
	PostgresUser                 string `envconfig:"postgres_user"`
	PostgresDB                   string `envconfig:"postgres_db"`
	MailgunApiKey                string `envconfig:"mg_public_api_key"`
	EmailFrom                    string `envconfig:"email_from"`
	BaseUrl                      string `envconfig:"base_url"`
	Env                          string `envconfig:"env"`
	PostgresPort                 int    `envconfig:"postgres_port"`
	PostgresPassword             string `envconfig:"postgres_password"`
	JWTSecret                    string `envconfig:"jwt_secret"`
	FacebookClientID             string `envconfig:"facebook_client_id"`
	FacebookClientSecret         string `envconfig:"facebook_client_secret"`
	FacebookRedirectURL          string `envconfig:"facebook_redirect_url"`
	MgDomain                     string `envconfig:"mg_domain"`
	Host                         string `envconfig:"host"`
	GoogleClientID               string `envconfig:"google_client_id"`
	GoogleClientSecret           string `envconfig:"google_client_secret"`
	GoogleRedirectURL            string `envconfig:"google_redirect_url"`
	GoogleApplicationCredentials string `envconfig:"google_application_credentials"`
}

func Load() (*Config, error) {
	env := os.Getenv("GIN_MODE")
	if env != "release" {
		if err := godotenv.Load("./.env"); err != nil {
			log.Printf("couldn't load env vars: %v", err)
		}
	}

	c := &Config{}
	err := envconfig.Process("upload", c)
	if err != nil {
		return nil, err
	}

	if c.Env == "prod" {
		// Fetch secret from Secret Manager
		secretValue, err := getSecret("projects/routepay/secrets/routepay-secret/versions/latest")
		if err != nil {
			return nil, err
		}

		// Set the sensitive value
		c.PostgresPassword = secretValue
	}

	return c, nil
}

func getSecret(secretName string) (string, error) {
	ctx := context.Background()

	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	secretVersion := "latest"
	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("%s/versions/%s", secretName, secretVersion),
	}

	result, err := client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		return "", err
	}

	return string(result.Payload.Data), nil
}
