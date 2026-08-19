package azs006

import (
	"github.com/example/provider/pluginsdk"
)

// azurerm_secretive: with the default flag settings a `Sensitive: true` resource property is
// still required of the data source, while a property carrying an //azignore:AZS006 directive
// is exempt without any flag.

func resourceSecretive() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"name": {Type: pluginsdk.TypeString, Required: true},
			"primary_access_key": {
				Type:      pluginsdk.TypeString,
				Sensitive: true,
			},
			"connection_string": { //azignore:AZS006
				Type:      pluginsdk.TypeString,
				Sensitive: true,
			},
			//azignore:AZS006
			"deployment_token": {
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func dataSourceSecretive() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"name": {Type: pluginsdk.TypeString, Required: true},
		},
	}
}
