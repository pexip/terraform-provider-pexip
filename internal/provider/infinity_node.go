package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var (
	_ resource.ResourceWithIdentity = (*infinityNodeResource)(nil)
)

type infinityNodeResource struct{}

func (r infinityNodeResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	//TODO implement me
	panic("implement me")
}

func (r infinityNodeResource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	//TODO implement me
	panic("implement me")
}

func (r infinityNodeResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	//TODO implement me
	panic("implement me")
}

func (r infinityNodeResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	//TODO implement me
	panic("implement me")
}

func (r infinityNodeResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	//TODO implement me
	panic("implement me")
}

func (r infinityNodeResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	//TODO implement me
	panic("implement me")
}

func (r infinityNodeResource) IdentitySchema(ctx context.Context, request resource.IdentitySchemaRequest, response *resource.IdentitySchemaResponse) {
	//TODO implement me
	panic("implement me")
}
