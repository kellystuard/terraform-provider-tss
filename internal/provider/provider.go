package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			DataSourcesMap: map[string]*schema.Resource{
				"tss_secret_field": dataSourceSecretField(),
			},
			// ResourcesMap: map[string]*schema.Resource{
			// 	"secret_resource": resourceSecret(),
			// },
			Schema: map[string]*schema.Schema{
				"username": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("TSS_USERNAME", nil),
				},
				"password": {
					Type:        schema.TypeString,
					Required:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("TSS_PASSWORD", nil),
				},
				"tenant": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("TSS_TENANT", nil),
				},
				"grant_type": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "password",
					DefaultFunc: schema.EnvDefaultFunc("TSS_GRANT_TYPE", nil),
				},
				"timeout": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "10s",
					DefaultFunc: schema.EnvDefaultFunc("TSS_TIMEOUT", nil),
				},
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

type apiClient struct {
	AccessToken string `json:"access_token"`
	BaseUrl     string
	Timeout     time.Duration
	UserAgent   string
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

		username := d.Get("username").(string)
		password := d.Get("password").(string)
		tenant := d.Get("tenant").(string)
		grant_type := d.Get("grant_type").(string)
		timeout, err := time.ParseDuration(d.Get("timeout").(string))

		config := &apiClient{
			BaseUrl:   fmt.Sprintf("https://%s.secretservercloud.com", tenant),
			Timeout:   timeout,
			UserAgent: p.UserAgent("terraform-provider-tss", version),
		}
		if err != nil {
			return config, diag.FromErr(err)
		}

		client := &http.Client{
			Timeout: config.Timeout,
		}

		var diags diag.Diagnostics

		body := strings.NewReader(fmt.Sprintf(
			"username=%s&password=%s&grant_type=%s",
			username, password, grant_type,
		))
		url := config.BaseUrl + "/oauth2/token"
		req, err := http.NewRequest("POST", url, body)
		if err != nil {
			return config, diag.FromErr(err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := client.Do(req)
		if err != nil {
			return config, diag.FromErr(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			return config, diag.Errorf("Oauth token response: (%d) %s", resp.StatusCode, body)
		}

		json.NewDecoder(resp.Body).Decode(config)
//return config, diag.Errorf("%s", config.AccessToken)

		return config, diags
	}
}
