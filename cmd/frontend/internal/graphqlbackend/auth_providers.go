package graphqlbackend

import (
	"context"

	"github.com/sourcegraph/sourcegraph/cmd/frontend/internal/backend"
	"github.com/sourcegraph/sourcegraph/pkg/conf"
	"github.com/sourcegraph/sourcegraph/schema"
)

func (r *siteResolver) AuthProviders(ctx context.Context) (*authProviderConnectionResolver, error) {
	// 🚨 SECURITY: Only site admins can list auth providers.
	if err := backend.CheckCurrentUserIsSiteAdmin(ctx); err != nil {
		return nil, err
	}
	return &authProviderConnectionResolver{
		authProviders: conf.AuthProviders(),
	}, nil
}

// authProviderConnectionResolver resolves a list of auth providers.
//
// 🚨 SECURITY: When instantiating an authProviderConnectionResolver value, the caller MUST check
// permissions.
type authProviderConnectionResolver struct {
	authProviders []schema.AuthProviders
}

func (r *authProviderConnectionResolver) Nodes(ctx context.Context) ([]*authProviderResolver, error) {
	// 🚨 SECURITY: Only site admins can list auth providers. This check is intentionally redundant
	// (with the check in (*siteResolver).AuthProviders), to reduce the likelihood of a bug causing
	// this information to leak.
	if err := backend.CheckCurrentUserIsSiteAdmin(ctx); err != nil {
		return nil, err
	}

	var rs []*authProviderResolver
	for _, authProvider := range r.authProviders {
		rs = append(rs, &authProviderResolver{authProvider: authProvider})
	}
	return rs, nil
}

func (r *authProviderConnectionResolver) TotalCount() int32   { return int32(len(r.authProviders)) }
func (r *authProviderConnectionResolver) PageInfo() *pageInfo { return &pageInfo{hasNextPage: false} }