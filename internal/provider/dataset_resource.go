// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"terraform-provider-zosmf/zosmf"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &DatasetResource{}
var _ resource.ResourceWithConfigure = &DatasetResource{}

func NewDatasetResource() resource.Resource {
	return &DatasetResource{}
}

// DatasetResource defines the resource implementation.
type DatasetResource struct {
	client *zosmf.Client
}

// DatasetResourceModel describes the resource data model.
type DatasetResourceModel struct {
	Name      types.String `tfsdk:"name"`
	Volser    types.String `tfsdk:"volser"`
	Unit      types.String `tfsdk:"unit"`
	Dsorg     types.String `tfsdk:"dsorg"`
	Alcunit   types.String `tfsdk:"alcunit"`
	Primary   types.Int64  `tfsdk:"primary"`
	Secondary types.Int64  `tfsdk:"secondary"`
	Dirblk    types.Int64  `tfsdk:"dirblk"`
	Avgblk    types.Int64  `tfsdk:"avgblk"`
	Recfm     types.String `tfsdk:"recfm"`
	Blksize   types.Int64  `tfsdk:"blksize"`
	Lrecl     types.Int64  `tfsdk:"lrecl"`
	Storclass types.String `tfsdk:"storclass"`
	Mgntclass types.String `tfsdk:"mgntclass"`
	Dataclass types.String `tfsdk:"dataclass"`
	Dsntype   types.String `tfsdk:"dsntype"`
}

func (r *DatasetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dataset"
}

func (r *DatasetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Dataset resource",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Dataset name",
				Optional:            false,
				Required:            true,
			},
			"volser": schema.StringAttribute{
				MarkdownDescription: "Volume",
				Optional:            true,
			},
			"unit": schema.StringAttribute{
				MarkdownDescription: "Device type",
				Optional:            true,
			},
			"dsorg": schema.StringAttribute{
				MarkdownDescription: "Data set organization",
				Optional:            true,
			},
			"alcunit": schema.StringAttribute{
				MarkdownDescription: "Unit of space allocation",
				Optional:            true,
			},
			"primary": schema.StringAttribute{
				MarkdownDescription: "Primary space allocation",
				Optional:            true,
			},
			"secondary": schema.StringAttribute{
				MarkdownDescription: "Secondary space allocation",
				Optional:            true,
			},
			"dirblk": schema.StringAttribute{
				MarkdownDescription: "Number of directory blocks",
				Optional:            true,
			},
			"avgblk": schema.StringAttribute{
				MarkdownDescription: "Average block",
				Optional:            true,
			},
			"recfm": schema.StringAttribute{
				MarkdownDescription: "Record format",
				Optional:            true,
			},
			"blksize": schema.StringAttribute{
				MarkdownDescription: "Block size",
				Optional:            true,
			},
			"lrecl": schema.StringAttribute{
				MarkdownDescription: "Record length",
				Optional:            true,
			},
			"storclass": schema.StringAttribute{
				MarkdownDescription: "Storage class",
				Optional:            true,
			},
			"mgntclass": schema.StringAttribute{
				MarkdownDescription: "Management class",
				Optional:            true,
			},
			"dataclass": schema.StringAttribute{
				MarkdownDescription: "Data class",
				Optional:            true,
			},
			"dsntype": schema.StringAttribute{
				MarkdownDescription: "Dataset type",
				Optional:            true,
			},
		},
	}
}

func (r *DatasetResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*zosmf.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *DatasetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DatasetResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	datasetattrib := zosmf.CreateDataset{
		Volser:    data.Volser.ValueString(),
		Unit:      data.Unit.ValueString(),
		Dsorg:     data.Dsorg.ValueString(),
		Alcunit:   data.Name.ValueString(),
		Primary:   int(data.Primary.ValueInt64()),
		Secondary: int(data.Secondary.ValueInt64()),
		Dirblk:    int(data.Dirblk.ValueInt64()),
		Avgblk:    int(data.Avgblk.ValueInt64()),
		Recfm:     data.Recfm.ValueString(),
		Blksize:   int(data.Blksize.ValueInt64()),
		Lrecl:     int(data.Lrecl.ValueInt64()),
		Storclass: data.Storclass.ValueString(),
		Mgntclass: data.Mgntclass.ValueString(),
		Dataclass: data.Dataclass.ValueString(),
		Dsntype:   data.Dsntype.ValueString(),
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	err := r.client.CreateDataset(data.Name.ValueString(), datasetattrib)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Dataset, got error: %s", err))
		return
	}

	// For the purposes of this Dataset code, hardcoding a response value to
	// save into the Terraform state.

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatasetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DatasetResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read Dataset, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatasetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DatasetResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update Dataset, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatasetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DatasetResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete Dataset, got error: %s", err))
	//     return
	// }
}

func (r *DatasetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
