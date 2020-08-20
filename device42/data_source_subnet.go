package device42

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSubnet() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSubnetRead,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"network": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"mask_bits": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"range_begin": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"range_end": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vrf_group_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"assigned": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"allocated": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"parent_subnet_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"subnet_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceSubnetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)
	var diags diag.Diagnostics

	name := d.Get("name").(string)

	subnet, err := c.GetSubnet(name)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("network", subnet.Subnets[0].Network)
	d.Set("mask_bits", subnet.Subnets[0].MaskBits)
	d.Set("range_begin", subnet.Subnets[0].RangeBegin)
	d.Set("range_end", subnet.Subnets[0].RangeEnd)
	d.Set("vrf_group_name", subnet.Subnets[0].VrfGroupName)
	d.Set("assigned", subnet.Subnets[0].Assigned)
	d.Set("allocated", subnet.Subnets[0].Allocated)
	d.Set("parent_subnet_id", subnet.Subnets[0].ParentSubnetID)
	d.Set("subnet_id", subnet.Subnets[0].SubnetID)

	sid := strconv.Itoa(subnet.Subnets[0].SubnetID)
	d.SetId(sid)

	return diags
}
