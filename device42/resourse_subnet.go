package device42

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSubnet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSubnetCreate,
		ReadContext:   resourceSubnetRead,
		UpdateContext: resourceSubnetUpdate,
		Delete:        resourceSubnetDelete,
		Schema: map[string]*schema.Schema{
			"last_updated": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"network": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"mask_bits": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"range_begin": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"range_end": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"vrf_group_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"assigned": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"allocated": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"parent_subnet_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"subnet_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"gateway": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func resourceSubnetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	pid := d.Get("parent_subnet_id").(int)
	mask := d.Get("mask_bits").(int)

	child, err := c.SetChild(pid, mask)
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.SetSubnet(child.Network, mask, name)
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.SetGateway(name)
	if err != nil {
		return diag.FromErr(err)
	}

	sid := strconv.Itoa(child.SubnetID)
	d.SetId(sid)
	resourceSubnetRead(ctx, d, m)

	return diags
}

func resourceSubnetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)
	var diags diag.Diagnostics

	name := d.Get("name").(string)

	subnet, err := c.GetSubnet(name)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("network", subnet.Subnets[0].Network)
	d.Set("range_begin", subnet.Subnets[0].RangeBegin)
	d.Set("range_end", subnet.Subnets[0].RangeEnd)
	d.Set("vrf_group_name", subnet.Subnets[0].VrfGroupName)
	d.Set("assigned", subnet.Subnets[0].Assigned)
	d.Set("allocated", subnet.Subnets[0].Allocated)
	d.Set("subnet_id", subnet.Subnets[0].SubnetID)
	d.Set("gateway", subnet.Subnets[0].Gateway)

	return diags
}

func resourceSubnetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	if d.HasChange("name") {
		name := d.Get("name").(string)
		net := d.Get("network").(string)
		mask := d.Get("mask_bits").(int)

		err := c.SetSubnet(net, mask, name)
		if err != nil {
			return diag.FromErr(err)
		}

		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	return resourceSubnetRead(ctx, d, m)
}

func resourceSubnetDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)

	sid := d.Id()

	err := c.DeleteSubnet(sid)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
