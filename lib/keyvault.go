package lib

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/apex/log"
	"github.com/imdario/mergo"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
)

func GetSecrets(keyVaultName string) (map[string]string, error) {
	ctx := log.WithField("vault_name", keyVaultName)

	keyVaultURL := fmt.Sprintf("https://%s.vault.azure.net/", keyVaultName)
	ctx.Debugf("getting credentials")

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to obtain Azure credentials: %v", err)
	}

	client, err := azsecrets.NewClient(keyVaultURL, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to client: %v", err)
	}

	secrets := make(map[string]string)

	//List secrets
	pager := client.ListPropertiesOfSecrets(nil)
	for pager.More() {
		page, err := pager.NextPage(context.TODO())
		if err != nil {
			return nil, fmt.Errorf("error while listing secrets: %s", err)
		}

		for _, v := range page.Secrets {
			secret, err := client.GetSecret(context.TODO(), *v.Name, nil)
			if err != nil {
				return nil, fmt.Errorf("can't retrieve '%s' from '%s': %s", *v.Name, keyVaultName, err)
			}

			keyValue := make(map[string]string)
			// if the content-type is application/json we'll get it as an object.
			// otherwise it'll be taken as a key->value where key is the secret name
			if contentType := v.ContentType; contentType != nil && *contentType == "application/json" {
				log.Debugf("JSON Secret Name: %s\tSecret Tags: %v", *v.ID, v.Tags)
				if err := json.Unmarshal([]byte(*secret.Value), &keyValue); err != nil {
					return nil, fmt.Errorf("Can't unmarshal json: %s", err)
				}
			} else {
				log.Debugf("PLAIN Secret Name: %s\tSecret Tags: %v", *v.ID, v.Tags)
				keyValue[*v.Name] = *secret.Value
			}

			// merge the resulting map with the collector one
			err = mergo.Merge(&secrets, &keyValue, mergo.WithOverride)
			if err != nil {
				return nil, fmt.Errorf("Can't merge maps: %s", err)
			}
		}
	}

	return secrets, nil
}
