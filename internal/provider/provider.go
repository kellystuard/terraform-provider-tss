package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
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
					Description: "Username passed to the OAuth token endpoint. Can be set through environment variable `TSS_USERNAME`.",
				},
				"password": {
					Type:        schema.TypeString,
					Required:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("TSS_PASSWORD", nil),
					Description: "Password passed to the OAuth token endpoint. Can be set through environment variable `TSS_USERNAME`.",
				},
				"tenant": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("TSS_TENANT", nil),
					Description: "Tenant used for all API communication. Set to the tenant portion of `https://tenant.secretservercloud.com/`. If using an on-premise installation, set to the full URI of the server (e.g. -- `https://my-server/SecretServer`). Can be set through environment variable `TSS_TENANT`.",
				},
				"grant_type": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "password",
					DefaultFunc: schema.EnvDefaultFunc("TSS_GRANT_TYPE", nil),
					Description: "Grant type passed to the OAuth token endpoint. Default is `password`. Can be set through environment variable `TSS_GRANT_TYPE`.",
				},
				"timeout": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "10s",
					DefaultFunc: schema.EnvDefaultFunc("TSS_TIMEOUT", nil),
					Description: "Timeout duration used for all API communication. Default is `10s` (10 seconds). Can be set through environment variable `TSS_TIMEOUT`.",
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

		log := hclog.Default()
		log.Trace("[TRACE] Configuring TSS provider")

		username := d.Get("username").(string)
		password := d.Get("password").(string)
		tenant := d.Get("tenant").(string)
		grant_type := d.Get("grant_type").(string)
		timeout, err := time.ParseDuration(d.Get("timeout").(string))

		var base_url string
		if strings.Contains(tenant, "//") {
			base_url = tenant
			log.Trace("Using on-premise instance", "tenant", tenant, "url", base_url)
		} else {
			base_url = fmt.Sprintf("https://%s.secretservercloud.com", tenant)
			log.Trace("Using cloud instance", "tenant", tenant, "url", base_url)
		}

		config := &apiClient{
			BaseUrl:   base_url,
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

		body := fmt.Sprintf(
			"username=%s&password=%s&grant_type=%s",
			username, password, grant_type,
		)
		url := config.BaseUrl + "/oauth2/token"

		log.Trace("Posting OAuth request", "url", url, "body", body, "timeout", config.Timeout)
		req, err := http.NewRequest("POST", url, strings.NewReader(body))
		if err != nil {
			return config, diag.FromErr(err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := client.Do(req)
		if err != nil {
			return config, diag.FromErr(err)
		}
		defer resp.Body.Close()

		log.Trace("Received OAuth response", "StatusCode", resp.StatusCode, "ContentLength", resp.ContentLength)
		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			return config, diag.Errorf("Oauth token response: (%d) %s", resp.StatusCode, body)
		}

		json.NewDecoder(resp.Body).Decode(config)
		log.Trace("OAuth response", "AccessToken", config.AccessToken)

		return config, diags
	}
}
