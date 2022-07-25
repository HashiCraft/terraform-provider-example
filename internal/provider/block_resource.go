package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-provider-scaffolding-framework/internal/client"
)

type blockResource struct {
	provider provider
}

type blockResourceData struct {
	X        types.Int64  `tfsdk:"x"`
	Y        types.Int64  `tfsdk:"y"`
	Z        types.Int64  `tfsdk:"z"`
	Material types.String `tfsdk:"material"`
	Id       types.String `tfsdk:"id"`
}

func (r *blockResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data blockResourceData

	// Convert the HCL config into the data struct
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	br := client.BlockRequest{
		X:        int(data.X.Value),
		Y:        int(data.Y.Value),
		Z:        int(data.Z.Value),
		Material: data.Material.Value,
	}

	err := r.provider.client.CreateBlock(br)
	if err != nil {
		resp.Diagnostics.AddError("unable to create block", err.Error())
		return
	}

	// block positions are unique, use this as the id
	data.Id = types.String{Value: fmt.Sprintf(
		"%d_%d_%d",
		data.X.Value,
		data.Y.Value,
		data.Z.Value),
	}

	tflog.Trace(ctx, "created a resource")

	// store the updated object in the state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *blockResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data blockResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	b, err := r.provider.client.GetBlock(int(data.X.Value), int(data.Y.Value), int(data.Z.Value))
	if err != nil {
		resp.Diagnostics.AddError("error fetching block", err.Error())
	}

	// block does not exist return
	if b == nil {
		return
	}

	data.Material = types.String{Value: b.Material}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *blockResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data blockResourceData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// blocks are immutable in minecraft to update a block we need to delete it then create it
	br := client.BlockRequest{
		X: int(data.X.Value),
		Y: int(data.Y.Value),
		Z: int(data.Z.Value),
	}

	err := r.provider.client.DeleteBlock(br)
	if err != nil {
		diags.AddError("unable to delete block", err.Error())
		return
	}

	// Material is currently the only changeable property
	br.Material = data.Material.Value

	err = r.provider.client.CreateBlock(br)
	if err != nil {
		diags.AddError("unable to delete block", err.Error())
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *blockResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data blockResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	br := client.BlockRequest{
		X: int(data.X.Value),
		Y: int(data.Y.Value),
		Z: int(data.Z.Value),
	}

	err := r.provider.client.DeleteBlock(br)
	if err != nil {
		diags.AddError("unable to delete block", err.Error())
		return
	}
}

type blockResourceType struct{}

func (t blockResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Minecraft Block",

		Attributes: map[string]tfsdk.Attribute{
			"x": {
				MarkdownDescription: "X position for the block",
				Optional:            true,
				Type:                types.Int64Type,
			},
			"y": {
				MarkdownDescription: "Y position for the block",
				Optional:            true,
				Type:                types.Int64Type,
			},
			"z": {
				MarkdownDescription: "Z position for the block",
				Optional:            true,
				Type:                types.Int64Type,
			},
			"material": {
				MarkdownDescription: "Material type for the block",
				Optional:            true,
				Type:                types.StringType,
			},
			"id": {
				Computed:            true,
				MarkdownDescription: "Identifier for the block",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
		},
	}, nil
}

func (t blockResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return &blockResource{
		provider: provider,
	}, diags
}
