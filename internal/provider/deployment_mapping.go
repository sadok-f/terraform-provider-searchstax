package provider

import (
	"context"

	searchstaxClient "terraform-provider-searchstax/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func deploymentCommonSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"uid":                              schema.StringAttribute{Computed: true},
		"name":                             schema.StringAttribute{Computed: true},
		"application":                      schema.StringAttribute{Computed: true},
		"application_version":              schema.StringAttribute{Computed: true},
		"tier":                             schema.StringAttribute{Computed: true},
		"http_endpoint":                    schema.StringAttribute{Computed: true},
		"status":                           schema.StringAttribute{Computed: true},
		"provision_state":                  schema.StringAttribute{Computed: true},
		"termination_lock":                 schema.BoolAttribute{Computed: true},
		"plan_type":                        schema.StringAttribute{Computed: true},
		"plan":                             schema.StringAttribute{Computed: true},
		"is_master_slave":                  schema.BoolAttribute{Computed: true},
		"vpc_type":                         schema.StringAttribute{Computed: true},
		"vpc_name":                         schema.StringAttribute{Computed: true},
		"region_id":                        schema.StringAttribute{Computed: true},
		"cloud_provider":                   schema.StringAttribute{Computed: true},
		"cloud_provider_id":                schema.StringAttribute{Computed: true},
		"num_additional_app_nodes":         schema.Int64Attribute{Computed: true},
		"deployment_type":                  schema.StringAttribute{Computed: true},
		"num_nodes_default":                schema.Int64Attribute{Computed: true},
		"num_zookeeper_nodes_default":      schema.Int64Attribute{Computed: true},
		"num_additional_zookeeper_nodes":   schema.Int64Attribute{Computed: true},
		"date_created":                     schema.StringAttribute{Computed: true},
		"servers":                          schema.ListAttribute{Computed: true, ElementType: types.StringType},
		"zookeeper_ensemble":               schema.StringAttribute{Computed: true},
		"tags":                             schema.ListAttribute{Computed: true, ElementType: types.StringType},
		"spec_jvm_heap_memory":             schema.StringAttribute{Computed: true},
		"spec_disk_space":                  schema.StringAttribute{Computed: true},
		"spec_physical_memory":             schema.StringAttribute{Computed: true},
		"backups_enabled":                  schema.BoolAttribute{Computed: true},
		"dr_enabled":                       schema.BoolAttribute{Computed: true},
		"sla_active":                       schema.BoolAttribute{Computed: true},
		"application_nodes_count":          schema.Int64Attribute{Computed: true},
		"subscription":                     schema.StringAttribute{Computed: true},
		"security_pack":                    schema.BoolAttribute{Computed: true},
		"desired_tier":                     schema.StringAttribute{Computed: true},
	}
}

type deploymentSpecificationsModel struct {
	JVMHeapMemory  types.String `tfsdk:"spec_jvm_heap_memory"`
	DiskSpace      types.String `tfsdk:"spec_disk_space"`
	PhysicalMemory types.String `tfsdk:"spec_physical_memory"`
}

type deploymentsModel struct {
	UID                         types.String                  `tfsdk:"uid"`
	Name                        types.String                  `tfsdk:"name"`
	Application                 types.String                  `tfsdk:"application"`
	ApplicationVersion          types.String                  `tfsdk:"application_version"`
	Tier                        types.String                  `tfsdk:"tier"`
	HttpEndpoint                types.String                  `tfsdk:"http_endpoint"`
	Status                      types.String                  `tfsdk:"status"`
	ProvisionState              types.String                  `tfsdk:"provision_state"`
	TerminationLock             types.Bool                    `tfsdk:"termination_lock"`
	Plan                        types.String                  `tfsdk:"plan"`
	PlanType                    types.String                  `tfsdk:"plan_type"`
	IsMasterSlave               types.Bool                    `tfsdk:"is_master_slave"`
	VpcType                     types.String                  `tfsdk:"vpc_type"`
	VpcName                     types.String                  `tfsdk:"vpc_name"`
	RegionId                    types.String                  `tfsdk:"region_id"`
	CloudProvider               types.String                  `tfsdk:"cloud_provider"`
	CloudProviderId             types.String                  `tfsdk:"cloud_provider_id"`
	NumAdditionalAppNodes       types.Int64                   `tfsdk:"num_additional_app_nodes"`
	DeploymentType              types.String                  `tfsdk:"deployment_type"`
	NumNodesDefault             types.Int64                   `tfsdk:"num_nodes_default"`
	NumZookeeperNodesDefault    types.Int64                   `tfsdk:"num_zookeeper_nodes_default"`
	NumAdditionalZookeeperNodes types.Int64                   `tfsdk:"num_additional_zookeeper_nodes"`
	DateCreated                 types.String                  `tfsdk:"date_created"`
	Servers                     types.List                    `tfsdk:"servers"`
	ZookeeperEnsemble           types.String                  `tfsdk:"zookeeper_ensemble"`
	Tags                        types.List                    `tfsdk:"tags"`
	SpecJVMHeapMemory           types.String                  `tfsdk:"spec_jvm_heap_memory"`
	SpecDiskSpace               types.String                  `tfsdk:"spec_disk_space"`
	SpecPhysicalMemory          types.String                  `tfsdk:"spec_physical_memory"`
	BackupsEnabled              types.Bool                    `tfsdk:"backups_enabled"`
	DrEnabled                   types.Bool                    `tfsdk:"dr_enabled"`
	SlaActive                   types.Bool                    `tfsdk:"sla_active"`
	ApplicationNodesCount       types.Int64                   `tfsdk:"application_nodes_count"`
	Subscription                types.String                  `tfsdk:"subscription"`
	SecurityPack                types.Bool                    `tfsdk:"security_pack"`
	DesiredTier                 types.String                  `tfsdk:"desired_tier"`
}

func mapDeploymentModel(ctx context.Context, deployment searchstaxClient.Deployment) (deploymentsModel, diag.Diagnostics) {
	servers, d := types.ListValueFrom(ctx, types.StringType, []string(deployment.Servers))
	if d.HasError() {
		return deploymentsModel{}, d
	}
	tags, d := types.ListValueFrom(ctx, types.StringType, []string(deployment.Tag))
	if d.HasError() {
		return deploymentsModel{}, d
	}
	specs := deploymentSpecificationsModel{
		JVMHeapMemory:  types.StringValue(deployment.Specifications.JVMHeapMemory),
		DiskSpace:      types.StringValue(deployment.Specifications.DiskSpace),
		PhysicalMemory: types.StringValue(deployment.Specifications.PhysicalMemory),
	}
	return deploymentsModel{
		UID:                         types.StringValue(deployment.UID),
		Name:                        types.StringValue(deployment.Name),
		Application:                 types.StringValue(deployment.Application),
		ApplicationVersion:          types.StringValue(deployment.ApplicationVersion),
		Tier:                        types.StringValue(deployment.Tier),
		Plan:                        types.StringValue(deployment.Plan),
		PlanType:                    types.StringValue(deployment.PlanType),
		HttpEndpoint:                types.StringValue(deployment.HttpEndpoint),
		Status:                      types.StringValue(deployment.Status),
		DateCreated:                 types.StringValue(deployment.DateCreated),
		CloudProvider:               types.StringValue(deployment.CloudProvider),
		CloudProviderId:             types.StringValue(deployment.CloudProviderId),
		DeploymentType:              types.StringValue(deployment.DeploymentType),
		ProvisionState:              types.StringValue(deployment.ProvisionState),
		RegionId:                    types.StringValue(deployment.RegionId),
		VpcType:                     types.StringValue(deployment.VpcType),
		VpcName:                     types.StringValue(deployment.VpcName),
		TerminationLock:             types.BoolValue(deployment.TerminationLock),
		IsMasterSlave:               types.BoolValue(deployment.IsMasterSlave),
		NumNodesDefault:             types.Int64Value(deployment.NumNodesDefault),
		NumAdditionalAppNodes:       types.Int64Value(deployment.NumAdditionalAppNodes),
		NumZookeeperNodesDefault:    types.Int64Value(deployment.NumZookeeperNodesDefault),
		NumAdditionalZookeeperNodes: types.Int64Value(deployment.NumAdditionalZookeeperNodes),
		Servers:                     servers,
		ZookeeperEnsemble:           types.StringValue(deployment.ZookeeperEnsemble),
		Tags:                        tags,
		SpecJVMHeapMemory:           specs.JVMHeapMemory,
		SpecDiskSpace:               specs.DiskSpace,
		SpecPhysicalMemory:          specs.PhysicalMemory,
		BackupsEnabled:              types.BoolValue(deployment.BackupsEnabled),
		DrEnabled:                   types.BoolValue(deployment.DrEnabled),
		SlaActive:                   types.BoolValue(deployment.SlaActive),
		ApplicationNodesCount:       types.Int64Value(deployment.ApplicationNodesCount),
		Subscription:                types.StringValue(deployment.Subscription),
		SecurityPack:                types.BoolValue(deployment.SecurityPack),
		DesiredTier:                 types.StringValue(deployment.DesiredTier),
	}, nil
}

func populateDeploymentResourceModel(ctx context.Context, target *deploymentModel, deployment searchstaxClient.Deployment) diag.Diagnostics {
	mapped, diags := mapDeploymentModel(ctx, deployment)
	if diags.HasError() {
		return diags
	}
	target.UID = mapped.UID
	target.Name = mapped.Name
	target.Application = mapped.Application
	target.ApplicationVersion = mapped.ApplicationVersion
	target.Tier = mapped.Tier
	target.HttpEndpoint = mapped.HttpEndpoint
	target.Status = mapped.Status
	target.ProvisionState = mapped.ProvisionState
	target.PlanType = mapped.PlanType
	target.Plan = mapped.Plan
	target.CloudProvider = mapped.CloudProvider
	target.CloudProviderId = mapped.CloudProviderId
	target.NumAdditionalAppNodes = mapped.NumAdditionalAppNodes
	target.NumNodesDefault = mapped.NumNodesDefault
	target.IsMasterSlave = mapped.IsMasterSlave
	target.VpcName = mapped.VpcName
	target.VpcType = mapped.VpcType
	target.RegionId = mapped.RegionId
	target.TerminationLock = mapped.TerminationLock
	target.DateCreated = mapped.DateCreated
	target.DeploymentType = mapped.DeploymentType
	target.NumZookeeperNodesDefault = mapped.NumZookeeperNodesDefault
	target.NumAdditionalZookeeperNodes = mapped.NumAdditionalZookeeperNodes
	target.Servers = mapped.Servers
	target.ZookeeperEnsemble = mapped.ZookeeperEnsemble
	target.Tags = mapped.Tags
	target.SpecJVMHeapMemory = mapped.SpecJVMHeapMemory
	target.SpecDiskSpace = mapped.SpecDiskSpace
	target.SpecPhysicalMemory = mapped.SpecPhysicalMemory
	target.BackupsEnabled = mapped.BackupsEnabled
	target.DrEnabled = mapped.DrEnabled
	target.SlaActive = mapped.SlaActive
	target.ApplicationNodesCount = mapped.ApplicationNodesCount
	target.Subscription = mapped.Subscription
	target.SecurityPack = mapped.SecurityPack
	target.DesiredTier = mapped.DesiredTier
	if deployment.PrivateVpc != 0 {
		target.PrivateVpc = types.Int64Value(deployment.PrivateVpc)
	}
	return nil
}
