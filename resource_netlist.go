package main

import (
	"log"
	"sort"

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
				Default:  "created by xakamai-tf",
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"items": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				ForceNew: false,
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

	sort.Strings(newListItems)

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
	c := m.(*edgegrid.Client)

	listID := d.Id()
	listType := d.Get("type").(string)

	listOpts := edgegrid.ListNetworkListsOptions{
		Extended:          false,
		IncludeDeprecated: false,
		TypeOflist:        listType,
		IncludeElements:   true,
	}

	akamaiNetworkList, resp, _ := c.NetworkLists.GetNetworkList(listID, listOpts)
	if resp.Response.StatusCode == 404 {
		log.Printf("[WARN] Akamai network list (%s) not found", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("type", listType)
	d.Set("description", akamaiNetworkList.Description)
	d.Set("name", akamaiNetworkList.Name)

	sort.Strings(akamaiNetworkList.List)
	d.Set("items", akamaiNetworkList.List)

	return nil
}

func resourceNetlistUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*edgegrid.Client)

	listID := d.Id()
	listType := d.Get("type").(string)

	listOpts := edgegrid.ListNetworkListsOptions{
		Extended:          false,
		IncludeDeprecated: false,
		TypeOflist:        listType,
		IncludeElements:   true,
	}

	existingNetList, _, err := c.NetworkLists.GetNetworkList(listID, listOpts)
	if err != nil {
		return err
	}

	if d.HasChange("description") {
		existingNetList.Description = d.Get("description").(string)
	}

	if d.HasChange("items") {
		modifiedItems := []string{}
		for _, item := range d.Get("items").([]interface{}) {
			modifiedItems = append(modifiedItems, item.(string))
		}

		sort.Strings(modifiedItems)
		existingNetList.List = modifiedItems
	}

	log.Printf("[DEBUG] update Akamai network list with ID %s", listID)

	_, _, error := c.NetworkLists.ModifyNetworkList(listID, *existingNetList)
	if error != nil {
		return err
	}

	return resourceNetlistRead(d, m)
}

func resourceNetlistDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceNetlistExists(d *schema.ResourceData, m interface{}) (bool, error) {
	c := m.(*edgegrid.Client)
	log.Printf("[DEBUG] read Akamai network list with ID %s", d.Id())

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
