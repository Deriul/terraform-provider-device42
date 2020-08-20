package device42

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"address": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("DEVICE42_ADDRESS", nil),
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DEVICE42_USERNAME", nil),
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("DEVICE42_PASSWORD", nil),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"device42_subnet": dataSourceSubnet(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"device42_subnet": resourceSubnet(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	addr := d.Get("address").(string)
	uname := d.Get("username").(string)
	pwd := d.Get("password").(string)

	c, err := NewClient(&addr, &uname, &pwd)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create device42 client",
			Detail:   "Unable to auth user for authenticated device42 client",
		})

		return nil, diags
	}

	return c, diags
}
