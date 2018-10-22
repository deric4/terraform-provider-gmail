package gmail

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/gmail/v1"
)

func resourceLabel() *schema.Resource {
	return &schema.Resource{
		Create: resourceLabelCreate,
		Read:   resourceLabelRead,
		Update: resourceLabelUpdate,
		Delete: resourceLabelDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"label_list_visibility": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"message_list_visibility": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceLabelCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	label, err := resourceLabelBuild(d, meta)
	if err != nil {
		return err
	}

	ctx, cancel := contextWithTimeout()
	defer cancel()
	labelAPI, err := config.gmail.Users.Labels.
		Create("me", label).
		Context(ctx).
		Do()
	if err != nil {
		fmt.Printf("%s", err)
		return err
	}

	d.SetId(labelAPI.Id)

	return resourceLabelRead(d, meta)
}

func resourceLabelRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	ctx, cancel := contextWithTimeout()
	defer cancel()

	label, err := config.gmail.Users.Labels.
		Get("me", d.Id()).
		Context(ctx).
		Do()
	if err != nil {
		return err
	}

	d.Set("name", label.Name)
	d.Set("id", label.Id)
	d.Set("label_list_visibility", label.LabelListVisibility)
	d.Set("message_list_visibility", label.MessageListVisibility)

	return nil
}

func resourceLabelUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	label, err := resourceLabelBuild(d, meta)
	if err != nil {
		return err
	}

	ctx, cancel := contextWithTimeout()
	defer cancel()

	labelAPI, err := config.gmail.Users.Labels.
		Update("me", d.Id(), label).
		Context(ctx).
		Do()
	if err != nil {
		return err
	}

	d.SetId(labelAPI.Id)

	return resourceLabelRead(d, meta)
}

func resourceLabelDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	label_id := d.Id()

	ctx, cancel := contextWithTimeout()
	defer cancel()
	err := config.gmail.Users.Labels.
		Delete("me", label_id).
		Context(ctx).
		Do()
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

// resourceBuildEvent is a shared helper function which builds an "label" struct
// from the schema. This is used by create and update.
func resourceLabelBuild(d *schema.ResourceData, meta interface{}) (*gmail.Label, error) {
	name := d.Get("name").(string)
	labelListVisibility := d.Get("label_list_visibility").(string)
	messageListVisibility := d.Get("message_list_visibility").(string)
	label_id := d.Id()

	var label gmail.Label

	label.Name = name
	label.LabelListVisibility = labelListVisibility
	label.MessageListVisibility = messageListVisibility
	label.Id = label_id

	return &label, nil
}
