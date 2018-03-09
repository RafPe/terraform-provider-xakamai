package main

import (
	"log"

	"github.com/RafPe/go-edgegrid"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceNetlist() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetlistCreate,
		Read:   resourceNetlistRead,
		Update: resourceNetlistUpdate,
		Delete: resourceNetlistDelete,
		Exists: resourceNetlistExists,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "created by tf-xakamai",
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"items": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceNetlistCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*edgegrid.Client)
	newListItems := []string{}
	for _, item := range d.Get("items").([]interface{}) {
		newListItems = append(newListItems, item.(string))
	}

	listName := d.Get("name").(string)
	listType := d.Get("type").(string)
	listDesc := d.Get("description").(string)

	newListOpts := edgegrid.CreateNetworkListOptions{
		Name:        listName,
		Type:        listType,
		Description: listDesc,
		List:        newListItems,
	}
	newList, _, err := c.NetworkLists.CreateNetworkList(newListOpts)
	if err != nil {
		return err
	}
	d.SetId(newList.UniqueID)
	return nil
}

func resourceNetlistRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceNetlistUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceNetlistDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceNetlistExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	c := meta.(*edgegrid.Client)
	log.Printf("[DEBUG] read gitlab group %s", d.Id())

	listID := d.Id()
	listType := d.Get("type").(string)

	listOpts := edgegrid.ListNetworkListsOptions{
		Extended:          false,
		IncludeDeprecated: false,
		TypeOflist:        listType,
		IncludeElements:   true,
	}

	exists, resp, err := c.NetworkLists.GetNetworkList(listID, listOpts)
	if resp.Response.StatusCode == 404 {
		// We ignore client returning 404 as list is just not found
		err = nil
	}

	return exists != nil, err
}
