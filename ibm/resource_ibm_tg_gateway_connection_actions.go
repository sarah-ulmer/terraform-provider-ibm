// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"time"

	"github.com/IBM/networking-go-sdk/transitgatewayapisv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	tgXacGatewayId     = "gateway"
	tgXacConnectionId  = "connection_id"
	tgConnectionAction = "action"
)

func resourceIBMTransitGatewayConnectionActions() *schema.Resource {
	return &schema.Resource{
		Create: resourceIBMTransitGatewayConnectionActionsCreate,
		Read:   resourceIBMTransitGatewayConnectionActionsRead,
		Delete: resourceIBMTransitGatewayConnectionActionsDelete,
		// Exists:   resourceIBMTransitGatewayConnectionActionExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			// Don't need timeout for every CRUD
		},

		Schema: map[string]*schema.Schema{
			tgGatewayId: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The Transit Gateway identifier",
			},
			tgConnectionId: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The Transit Gateway Connection identifier",
			},
			tgConnectionAction: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: InvokeValidator("ibm_tg_connection_actions", tgConnectionAction),
				Description:  "The Transit Gateway Connection cross account action",
			},
		},
	}
}

func resourceIBMTransitGatewayConnectionActionsValidator() *ResourceValidator {

	validateSchema := make([]ValidateSchema, 0)
	actions := "approve, reject"
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 tgConnectionAction,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              actions})

	ibmTransitGatewayConnectionActionsResourceValidator := ResourceValidator{ResourceName: "ibm_tg_connection_actions", Schema: validateSchema}

	return &ibmTransitGatewayConnectionActionsResourceValidator
}

func resourceIBMTransitGatewayConnectionActionsCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := transitgatewayClient(meta)
	if err != nil {
		return err
	}
	createTransitGatewayConnectionActionsOptions := &transitgatewayapisv1.CreateTransitGatewayConnectionActionsOptions{}
	gatewayId := d.Get(tgXacGatewayId).(string)
	createTransitGatewayConnectionActionsOptions.SetTransitGatewayID(gatewayId)
	connectionId := d.Get(tgXacConnectionId).(string)
	createTransitGatewayConnectionActionsOptions.SetID(connectionId)
	action := d.Get(tgConnectionAction).(string)
	createTransitGatewayConnectionActionsOptions.SetAction(action)

	response, err := client.CreateTransitGatewayConnectionActions(createTransitGatewayConnectionActionsOptions)
	if err != nil {
		return fmt.Errorf("XAC connection action err %s\n%s", err, response)
	}

	// create sometimes returns read
	// delete returns nil
	d.SetId("")
	return nil
}

func resourceIBMTransitGatewayConnectionActionsRead(d *schema.ResourceData, meta interface{}) error {

	client, err := transitgatewayClient(meta)
	if err != nil {
		return err
	}
	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}

	gatewayId := parts[0]
	ID := parts[1]

	getTransitGatewayConnectionOptions := &transitgatewayapisv1.GetTransitGatewayConnectionOptions{}
	getTransitGatewayConnectionOptions.SetTransitGatewayID(gatewayId)
	getTransitGatewayConnectionOptions.SetID(ID)
	instance, response, err := client.GetTransitGatewayConnection(getTransitGatewayConnectionOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Getting Transit Gateway Connection (%s): %s\n%s", ID, err, response)
	}

	if instance.RequestStatus != nil {
		d.Set(tgRequestStatus, *instance.RequestStatus)
	}
	d.Set(tgConnectionId, *instance.ID)
	d.Set(tgGatewayId, gatewayId)
	getTransitGatewayOptions := &transitgatewayapisv1.GetTransitGatewayOptions{
		ID: &gatewayId,
	}
	tgw, response, err := client.GetTransitGateway(getTransitGatewayOptions)
	if err != nil {
		return fmt.Errorf("Error Getting Transit Gateway : %s\n%s", err, response)
	}
	d.Set(RelatedCRN, *tgw.Crn)

	return nil
}

func resourceIBMTransitGatewayConnectionActionsDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}
