package mysql

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceRole() *schema.Resource {
	return &schema.Resource{
		Create: CreateRole,
		Read:   ReadRole,
		Delete: DeleteRole,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func CreateRole(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*MySQLConfiguration).Db

	roleName := d.Get("name").(string)

	stmtSQL := fmt.Sprintf("CREATE ROLE '%s'", roleName)
	log.Printf("[DEBUG] SQL: %s", stmtSQL)

	_, err := db.Exec(stmtSQL)
	if err != nil {
		return fmt.Errorf("error creating role: %s", err)
	}

	d.SetId(roleName)

	return nil
}

func ReadRole(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*MySQLConfiguration).Db

	stmtSQL := fmt.Sprintf("SHOW GRANTS FOR '%s'", d.Id())
	log.Printf("[DEBUG] SQL: %s", stmtSQL)

	_, err := db.Exec(stmtSQL)
	if err != nil {
		log.Printf("[WARN] Role (%s) not found; removing from state", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("name", d.Id())

	return nil
}

func DeleteRole(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*MySQLConfiguration).Db

	stmtSQL := fmt.Sprintf("DROP ROLE '%s'", d.Get("name").(string))
	log.Printf("[DEBUG] SQL: %s", stmtSQL)

	_, err := db.Exec(stmtSQL)
	if err != nil {
		return err
	}

	return nil
}
