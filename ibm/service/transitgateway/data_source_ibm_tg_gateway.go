// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package transitgateway

import (
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/networking-go-sdk/transitgatewayapisv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	tgConnName             = "name"
	tgConnections          = "connections"
	ID                     = "id"
	tgBaseNetworkType = "base_network_type"
	tgPrefixFiltersDefault = "prefix_filters_default"
)

func DataSourceIBMTransitGateway() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMTransitGatewayRead,

		Schema: map[string]*schema.Schema{
			tgName: {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The Transit Gateway identifier",
				ValidateFunc: validate.InvokeValidator("ibm_tg_gateway", tgName),
			},
			tgCrn: {
				Type:     schema.TypeString,
				Computed: true,
			},
			tgLocation: {
				Type:     schema.TypeString,
				Computed: true,
			},
			tgCreatedAt: {
				Type:     schema.TypeString,
				Computed: true,
			},
			tgGlobal: {
				Type:     schema.TypeBool,
				Computed: true,
			},
			tgStatus: {
				Type:     schema.TypeString,
				Computed: true,
			},
			tgUpdatedAt: {
				Type:     schema.TypeString,
				Computed: true,
			},
			tgResourceGroup: {
				Type:     schema.TypeString,
				Computed: true,
			},
			tgConnections: {
				Type:        schema.TypeList,
				Description: "Collection of transit gateway connections",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						ID: {
							Type:     schema.TypeString,
							Computed: true,
						},
						tgNetworkAccountID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the account which owns the network that is being connected. Generally only used if the network is in a different account than the gateway.",
						},
						tgNetworkId: {
							Type:     schema.TypeString,
							Computed: true,
						},
						tgConnName: {
							Type:     schema.TypeString,
							Computed: true,
						},
						tgNetworkType: {
							Type:     schema.TypeString,
							Computed: true,
						},
						tgBaseConnectionId: {
							Type:     schema.TypeString,
							Computed: true,
						},
						tgBaseNetworkType: {
							Type:     schema.TypeString,
							Computed: true,
						},
						tgLocalBgpAsn: {
							Type:     schema.TypeInt,
							Computed: true,
						},
						tgLocalGatewayIp: {
							Type:     schema.TypeString,
							Computed: true,
						},
						tgLocalTunnelIp: {
							Type:     schema.TypeString,
							Computed: true,
						},
						tgPrefixFilters: {
							Type:        schema.TypeList,
							Description: "Collection of prefix filters",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									tgID: {
										Type:     schema.TypeString,
										Computed: true,
									},
									tgAction: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Whether to permit or deny the prefix filter",
									},
									tgBefore: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Identifier of prefix filter that handles ordering",
									},
									tgCreatedAt: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The date and time the prefix filter was created",
									},
									tgGe: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "IP Prefix GE",
									},
									tgLe: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "IP Prefix LE",
									},
									tgPrefix: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "IP Prefix",
									},
									tgUpdatedAt: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The date and time that this prefix filter was last updated",
									},
								},
							},
						},
						tgPrefixFiltersDefault: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Default setting of permit or deny which applies to any routes that don't match a specified filter",
						},
						tgRemoteBgpAsn: {
							Type:     schema.TypeInt,
							Computed: true,
						},
						tgRemoteGatewayIp: {
							Type:     schema.TypeString,
							Computed: true,
						},
						tgRemoteTunnelIp: {
							Type:     schema.TypeString,
							Computed: true,
						},
						tgZone: {
							Type:     schema.TypeString,
							Computed: true,
						},
						tgMtu: {
							Type:     schema.TypeInt,
							Computed: true,
						},
						tgConectionCreatedAt: {
							Type:     schema.TypeString,
							Computed: true,
						},
						tgConnectionStatus: {
							Type:     schema.TypeString,
							Computed: true,
						},
						tgUpdatedAt: {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMTransitGatewayRead(d *schema.ResourceData, meta interface{}) error {

	client, err := transitgatewayClient(meta)
	if err != nil {
		return err
	}

	listTransitGatewaysOptionsModel := &transitgatewayapisv1.ListTransitGatewaysOptions{}
	listTransitGateways, response, err := client.ListTransitGateways(listTransitGatewaysOptionsModel)
	if err != nil {
		return fmt.Errorf("[ERROR] Error while listing transit gateways %s\n%s", err, response)
	}

	gwName := d.Get(tgName).(string)
	var foundGateway bool
	for _, tgw := range listTransitGateways.TransitGateways {

		if *tgw.Name == gwName {
			d.SetId(*tgw.ID)
			d.Set(tgCrn, tgw.Crn)
			d.Set(tgName, tgw.Name)
			d.Set(tgLocation, tgw.Location)
			d.Set(tgCreatedAt, tgw.CreatedAt.String())

			if tgw.UpdatedAt != nil {
				d.Set(tgUpdatedAt, tgw.UpdatedAt.String())
			}
			d.Set(tgGlobal, tgw.Global)
			d.Set(tgStatus, tgw.Status)

			if tgw.ResourceGroup != nil {
				rg := tgw.ResourceGroup
				d.Set(tgResourceGroup, *rg.ID)
			}
			foundGateway = true
		}
	}

	if !foundGateway {
		return fmt.Errorf("[ERROR] Couldn't find any gateway with the specified name: (%s)", gwName)
	}

	return dataSourceIBMTransitGatewayConnectionsRead(d, meta)

}
func dataSourceIBMTransitGatewayConnectionsRead(d *schema.ResourceData, meta interface{}) error {

	client, err := transitgatewayClient(meta)
	if err != nil {
		return err
	}
	listTransitGatewayConnectionsOptions := &transitgatewayapisv1.ListTransitGatewayConnectionsOptions{}
	tgGatewayId := d.Id()
	log.Println("tgGatewayId: ", tgGatewayId)

	listTransitGatewayConnectionsOptions.SetTransitGatewayID(tgGatewayId)
	listTGConnections, response, err := client.ListTransitGatewayConnections(listTransitGatewayConnectionsOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error while listing transit gateway connections %s\n%s", err, response)
	}
	connections := make([]map[string]interface{}, 0)

	for _, instance := range listTGConnections.Connections {
		tgConn := map[string]interface{}{}

		if instance.ID != nil {
			tgConn[ID] = *instance.ID
		}
		if instance.Name != nil {
			tgConn[tgConnName] = *instance.Name
		}
		if instance.NetworkType != nil {
			tgConn[tgNetworkType] = *instance.NetworkType
		}

		if instance.NetworkID != nil {
			tgConn[tgNetworkId] = *instance.NetworkID
		}
		if instance.NetworkAccountID != nil {
			tgConn[tgNetworkAccountID] = *instance.NetworkAccountID
		}

		if instance.BaseConnectionID != nil {
			tgConn[tgBaseConnectionId] = *instance.BaseConnectionID
		}

		if instance.BaseNetworkType != nil {
			tgConn[tgBaseNetworkType] = *instance.BaseNetworkType
		}

		if instance.LocalBgpAsn != nil {
			tgConn[tgLocalBgpAsn] = *instance.LocalBgpAsn
		}

		if instance.LocalGatewayIp != nil {
			tgConn[tgLocalGatewayIp] = *instance.LocalGatewayIp
		}

		if instance.LocalTunnelIp != nil {
			tgConn[tgLocalTunnelIp] = *instance.LocalTunnelIp
		}

		if instance.RemoteBgpAsn != nil {
			tgConn[tgRemoteBgpAsn] = *instance.RemoteBgpAsn
		}

		if instance.RemoteGatewayIp != nil {
			tgConn[tgRemoteGatewayIp] = *instance.RemoteGatewayIp
		}

		if instance.RemoteTunnelIp != nil {
			tgConn[tgRemoteTunnelIp] = *instance.RemoteTunnelIp
		}

		if instance.Zone != nil {
			tgConn[tgZone] = *instance.Zone.Name
		}

		if instance.Mtu != nil {
			tgConn[tgMtu] = *instance.Mtu
		}

		if instance.CreatedAt != nil {
			tgConn[tgConectionCreatedAt] = instance.CreatedAt.String()

		}
		if instance.UpdatedAt != nil {
			tgConn[tgUpdatedAt] = instance.UpdatedAt.String()

		}
		if instance.Status != nil {
			tgConn[tgConnectionStatus] = *instance.Status
		}

		prefixFiltersCollection := make([]map[string]interface{}, 0)
		for _, prefixFilter := range instance.PrefixFilters {
			tgPrefixFilter := map[string]interface{}{}
			tgPrefixFilter[tgID] = prefixFilter.ID
			tgPrefixFilter[tgAction] = prefixFilter.Action
			tgPrefixFilter[tgCreatedAt] = prefixFilter.CreatedAt.String()
			tgPrefixFilter[tgPrefix] = prefixFilter.Prefix

			if prefixFilter.UpdatedAt != nil {
				tgPrefixFilter[tgUpdatedAt] = prefixFilter.UpdatedAt.String()
			}
			if prefixFilter.Before != nil {
				tgPrefixFilter[tgBefore] = prefixFilter.Before
			}
			if prefixFilter.Ge != nil {
				tgPrefixFilter[tgGe] = prefixFilter.Ge
			}
			if prefixFilter.Le != nil {
				tgPrefixFilter[tgLe] = prefixFilter.Le
			}

			prefixFiltersCollection = append(prefixFiltersCollection, tgPrefixFilter)
		}
		tgConn[tgPrefixFilters] = prefixFiltersCollection

		if instance.PrefixFiltersDefault != nil {
			tgConn[tgPrefixFiltersDefault] = instance.PrefixFiltersDefault
		}

		connections = append(connections, tgConn)

	}
	d.Set(tgConnections, connections)

	return nil

}
