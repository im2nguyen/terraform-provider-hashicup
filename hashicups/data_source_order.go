package hashicups

import (
	"strconv"
	"context"

	hc "github.com/hashicorp-demoapp/hashicups-client-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func dataSourceOrder() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOrderRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"items": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"coffee": &schema.Schema{
							Type:     schema.TypeList,
							MaxItems: 1,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:     schema.TypeInt,
										Computed: true,
									},
									"name": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"teaser": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"description": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"price": &schema.Schema{
										Type:     schema.TypeInt,
										Computed: true,
									},
									"image": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"quantity": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceOrderRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*hc.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	orderID := strconv.Itoa(d.Get("id").(int))

	order, err := c.GetOrder(orderID)
	if err != nil {
		return diag.FromErr(err)
	}

	orderItems := flattenOrderItems(&order.Items)
	if err := d.Set("items", orderItems); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(orderID)

	return diags
}

func flattenOrderItems(orderItems *[]hc.OrderItem) []interface{} {
	if orderItems != nil {
		ois := make([]interface{}, len(*orderItems), len(*orderItems))

		for i, orderItem := range *orderItems {
			oi := make(map[string]interface{})

			oi["coffee"] = flattenCoffee(orderItem.Coffee)
			oi["quantity"] = orderItem.Quantity

			ois[i] = oi
		}

		return ois
	}

	return make([]interface{}, 0)
}

func flattenCoffee(coffee hc.Coffee) []interface{} {
	c := make(map[string]interface{})
	c["id"] = coffee.ID
	c["name"] = coffee.Name
	c["teaser"] = coffee.Teaser
	c["description"] = coffee.Description
	c["price"] = coffee.Price
	c["image"] = coffee.Image

	return []interface{}{c}
}
