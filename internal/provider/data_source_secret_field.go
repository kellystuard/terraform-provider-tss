package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSecretField() *schema.Resource {
	return &schema.Resource{
		Description: "Thycotic Secret Server data source in the Terraform provider for secret fields.",

		ReadContext: dataSourceSecretFieldRead,

		Schema: map[string]*schema.Schema{
			"number": {
				Description: "Secret ID",
				Required:    true,
				Type:        schema.TypeInt,
			},
			"slug": {
				Description: "Secret field name",
				Required:    true,
				Type:        schema.TypeString,
			},
			"value": {
				Computed:    true,
				Description: "Secret field value",
				Sensitive:   true,
				Type:        schema.TypeString,
			},
		},
	}
}

func dataSourceSecretFieldRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	log := hclog.Default()
	log.Trace("[TRACE] Reading secret field")

	// use the meta value to retrieve your client from the provider configure method
	config := meta.(*apiClient)
	client := &http.Client{Timeout: config.Timeout}

	var diags diag.Diagnostics

	number := d.Get("number").(int)
	slug := d.Get("slug")
	url := fmt.Sprintf("%s/api/v1/secrets/%d/fields/%s", config.BaseUrl, number, slug)

	log.Trace("[TRACE] Reading secret field from API", "url", url, "access token", config.AccessToken, "timeout", config.Timeout)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil{
		return diag.FromErr(err)
	}

	req.Header.Add("Authorization", "Bearer "+config.AccessToken)

	resp, err := client.Do(req)
	if err != nil{
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200{
		return diag.Errorf("Secret field response: %d", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	var result string
	if decoder.Decode(&result) != nil{
		return diag.FromErr(err)
	}

	if err := d.Set("value", result); err != nil{
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d/%s", number, slug))

	return diags
}
