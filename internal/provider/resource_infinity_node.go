package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/pexip/terraform-provider-pexip/internal/helpers"
	"github.com/rs/zerolog/log"
	"time"
)

func infinityNodeResourceQuery() *schema.Resource {
	return &schema.Resource{
		CreateContext: infinityNodeResourceQueryCreate,
		ReadContext:   infinityNodeResourceQueryRead,
		UpdateContext: infinityNodeResourceQueryUpdate,
		DeleteContext: infinityNodeResourceQueryDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Second),
			Update: schema.DefaultTimeout(10 * time.Second),
			Delete: schema.DefaultTimeout(10 * time.Second),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"inventory": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"group": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"variables": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.NoZeroValues,
				},
			},
		},
	}
}

func infinityNodeResourceQueryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(providerConfiguration)
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	name := helpers.ResourceToString(d, "name")
	groupID := helpers.ResourceToString(d, "group")
	inventoryRef := helpers.ResourceToString(d, "inventory")
	variables := helpers.ResourceToInterfaceMap(d, "variables")

	conf.Mutex.Lock()
	i, err := inventory.Load(conf.Path, inventoryRef)
	if err != nil {
		return diag.Errorf("failed to load inventory '%s': %s", inventoryRef, err.Error())
	}
	db, err := i.GetAndLoadDatabase()
	if err != nil {
		log.Error().Err(err).Msg("failed to load database")
		return diag.Errorf("failed to load database '%s': %s", inventoryRef, err.Error())
	}

	g := db.Group(groupID)
	if g == nil {
		return diag.Errorf("unable to find group '%s'", groupID)
	}

	h := database.NewHost(name, variables)
	g.UpdateEntity(h)
	db.UpdateGroup(*g)

	// Save and export database
	if err := commitAndExport(db, i.GetInventoryPath()); err != nil {
		return diag.FromErr(err)
	}
	conf.Mutex.Unlock()

	d.SetId(h.GetID())
	d.MarkNewResource()
	return infinityNodeResourceQueryRead(ctx, d, meta)
}

func infinityNodeResourceQueryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conf := meta.(providerConfiguration)
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	inventoryRef := helpers.ResourceToString(d, "inventory")

	conf.Mutex.Lock()
	i, err := inventory.Load(conf.Path, inventoryRef)
	if err != nil {
		return diag.Errorf("failed to load inventory '%s': %s", inventoryRef, err.Error())
	}
	db, err := i.GetAndLoadDatabase()
	conf.Mutex.Unlock()
	if err != nil {
		log.Error().Err(err).Msg("failed to load database")
		return diag.Errorf("failed to load database '%s': %s", inventoryRef, err.Error())
	}

	id := d.Id()
	g, entry, err := db.FindEntryByID(id)
	if err != nil {
		return diag.Errorf("unable to find entry '%s': %s", id, err.Error())
	}

	_ = d.Set("name", entry.GetName())
	_ = d.Set("group", g.GetID())

	h, ok := entry.(*database.Host)
	if ok {
		_ = d.Set("variables", h.GetVariables())
	}
	return diags
}

func infinityNodeResourceQueryUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(providerConfiguration)
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	name := helpers.ResourceToString(d, "name")
	groupID := helpers.ResourceToString(d, "group")
	inventoryRef := helpers.ResourceToString(d, "inventory")
	variables := helpers.ResourceToInterfaceMap(d, "variables")

	conf.Mutex.Lock()
	i, err := inventory.Load(conf.Path, inventoryRef)
	if err != nil {
		return diag.Errorf("failed to load inventory '%s': %s", inventoryRef, err.Error())
	}
	db, err := i.GetAndLoadDatabase()
	if err != nil {
		log.Error().Err(err).Msg("failed to load database")
		return diag.Errorf("failed to load database '%s': %s", inventoryRef, err.Error())
	}

	g, entry, err := db.FindEntryByID(d.Id())
	if err != nil {
		return diag.Errorf("unable to find entry '%s': %s", d.Id(), err.Error())
	}

	// check if name has changed
	if d.HasChange("name") {
		entry.SetName(name)
		db.UpdateGroup(*g)
	}

	// check if group has changed
	if d.HasChange("group") {
		// remove host from old group
		if err := g.RemoveEntity(entry); err != nil {
			return diag.Errorf("failed remove entry from group '%s': %s", g.GetID(), err.Error())
		}
		db.UpdateGroup(*g)

		// load new group
		ng := db.Group(groupID)
		if ng == nil {
			return diag.Errorf("failed to locate group '%s': %s", groupID, err.Error())
		}

		// update name and add entity to new group
		ng.UpdateEntity(entry)
	}

	if d.HasChange("variables") {
		h, ok := entry.(*database.Host)
		if ok {
			for k := range variables {
				h.SetVariable(k, variables[k])
			}
		}
		db.UpdateGroup(*g)
	}

	if d.HasChanges("name", "group", "variables") {
		// Save and export database
		if err := commitAndExport(db, i.GetInventoryPath()); err != nil {
			return diag.FromErr(err)
		}
	}

	conf.Mutex.Unlock()

	return infinityNodeResourceQueryRead(ctx, d, meta)
}

func infinityNodeResourceQueryDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conf := meta.(providerConfiguration)
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	inventoryRef := helpers.ResourceToString(d, "inventory")

	conf.Mutex.Lock()
	i, err := inventory.Load(conf.Path, inventoryRef)
	if err != nil {
		return diag.Errorf("failed to load inventory '%s': %s", inventoryRef, err.Error())
	}
	db, err := i.GetAndLoadDatabase()
	if err != nil {
		log.Error().Err(err).Msg("failed to load database")
		return diag.Errorf("failed to load database '%s': %s", inventoryRef, err.Error())
	}

	id := d.Id()
	g, entry, err := db.FindEntryByID(id)
	if err != nil {
		log.Error().Err(err).Msg("cannot find host so unable to remove, but continuing anyway")
	} else {
		// only remove host from group if we actually find it there. if we dont find it, then everything is ok and we
		// can skip the removing it.

		// remove entry from group
		if err := g.RemoveEntity(entry); err != nil {
			return diag.Errorf("unable to remove entry from group with id: %s", err.Error())
		}

		// update group
		db.UpdateGroup(*g)
	}

	// Save and export database
	if err := commitAndExport(db, i.GetInventoryPath()); err != nil {
		return diag.FromErr(err)
	}
	conf.Mutex.Unlock()

	return diags
}
