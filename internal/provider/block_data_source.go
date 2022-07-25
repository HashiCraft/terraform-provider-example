package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type blockDataSource struct {
	provider provider
}

func (d *blockDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	// reuse the same type from the resource
	var data blockResourceData

	tflog.Info(ctx, "Fetch data for datasource")

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// decode the id into separarate parts
	parts := strings.Split(data.Id.Value, "_")
	x, _ := strconv.Atoi(parts[0])
	y, _ := strconv.Atoi(parts[1])
	z, _ := strconv.Atoi(parts[2])

	b, err := d.provider.client.GetBlock(x, y, z)

	if b == nil {
		diags.AddError("Unable to find block", fmt.Sprintf("A block does not exist at %d,%d,%d", x, y, z))
		return
	}

	if err != nil {
		diags.AddError("Error fetching block", err.Error())
		return
	}

	// set the properties
	data.X = types.Int64{Value: int64(b.X)}
	data.Y = types.Int64{Value: int64(b.Y)}
	data.Z = types.Int64{Value: int64(b.Z)}
	data.Material = types.String{Value: b.Material}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

type blockDataSourceType struct{}

func (t *blockDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Minecraft Block",

		Attributes: map[string]tfsdk.Attribute{
			"x": {
				MarkdownDescription: "X position for the block",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"y": {
				MarkdownDescription: "Y position for the block",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"z": {
				MarkdownDescription: "Z position for the block",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"material": {
				MarkdownDescription: "Material type for the block",
				Computed:            true,
				Type:                types.StringType,
			},
			"id": {
				Required:            true,
				MarkdownDescription: "Identifier for the block",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
		},
	}, nil
}

func (t *blockDataSourceType) NewDataSource(ctx context.Context, in tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return &blockDataSource{
		provider: provider,
	}, diags
}
