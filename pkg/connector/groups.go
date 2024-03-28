package connector

import (
	"context"
	"fmt"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	ent "github.com/conductorone/baton-sdk/pkg/types/entitlement"
	grant "github.com/conductorone/baton-sdk/pkg/types/grant"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/conductorone/baton-verkada/pkg/verkada"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

const memberRole = "member"

type groupBuilder struct {
	resourceType *v2.ResourceType
	client       *verkada.Client
}

func (g *groupBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return g.resourceType
}

// Create a new connector resource for a Verkada group.
func groupResource(group *verkada.Group) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"group_name": group.Name,
		"group_id":   group.GroupID,
	}

	groupTraitOptions := []rs.GroupTraitOption{
		rs.WithGroupProfile(profile),
	}

	ret, err := rs.NewGroupResource(
		group.Name,
		groupResourceType,
		group.GroupID,
		groupTraitOptions,
	)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// List returns all the access groups from the database as resource objects.
func (g *groupBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	groups, err := g.client.ListAccessGroups(ctx)
	if err != nil {
		return nil, "", nil, fmt.Errorf("error listing access groups: %w", err)
	}

	var rv []*v2.Resource
	for _, group := range groups {
		groupCopy := group
		tr, err := groupResource(&groupCopy)
		if err != nil {
			return nil, "", nil, fmt.Errorf("error creating access group resource: %w", err)
		}
		rv = append(rv, tr)
	}

	return rv, "", nil, nil
}

func (g *groupBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	var rv []*v2.Entitlement
	options := []ent.EntitlementOption{
		ent.WithGrantableTo(userResourceType),
		ent.WithDisplayName(fmt.Sprintf("%s Group %s", resource.DisplayName, memberRole)),
		ent.WithDescription(fmt.Sprintf("%s of %s Verkada group", memberRole, resource.DisplayName)),
	}

	rv = append(rv, ent.NewAssignmentEntitlement(resource, memberRole, options...))

	return rv, "", nil, nil
}

func (g *groupBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)
	users, err := g.client.ListUsers(ctx)
	if err != nil {
		return nil, "", nil, fmt.Errorf("error getting users: %w", err)
	}

	var rv []*v2.Grant
	for _, user := range users {
		userCopy := user
		accessInfo, err := g.client.GetUserAccessInformation(ctx, user.UserID)
		if err != nil {
			l.Warn(
				"baton-verkada: error fetching user information, skipping user grant for group membership",
				zap.String("user_id", user.UserID),
				zap.String("group_id", resource.Id.Resource),
			)
			continue
		}

		if groupContainsUser(resource.Id.Resource, accessInfo.AccessGroups) {
			ur, err := userResource(&userCopy)
			if err != nil {
				return nil, "", nil, fmt.Errorf("error creating user resource for group %s: %w", resource.Id.Resource, err)
			}
			gr := grant.NewGrant(resource, memberRole, ur.Id)
			rv = append(rv, gr)
		}
	}
	return rv, "", nil, nil
}

func (g *groupBuilder) Grant(ctx context.Context, principal *v2.Resource, entitlement *v2.Entitlement) (annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)

	if principal.Id.ResourceType != userResourceType.Id {
		l.Warn(
			"baton-verkada: only users can be granted group membership",
			zap.String("principal_type", principal.Id.ResourceType),
			zap.String("principal_id", principal.Id.Resource),
		)
		return nil, fmt.Errorf("baton-verkada: only users can be granted group membership")
	}

	err := g.client.AddUserToGroup(ctx, entitlement.Resource.Id.Resource, principal.Id.Resource)
	if err != nil {
		return nil, fmt.Errorf("baton-verkada: failed to add user to group: %w", err)
	}

	return nil, nil
}

func (g *groupBuilder) Revoke(ctx context.Context, grant *v2.Grant) (annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)
	principal := grant.Principal
	entitlement := grant.Entitlement

	if principal.Id.ResourceType != userResourceType.Id {
		l.Warn(
			"baton-verkada: only users can have group membership revoked",
			zap.String("principal_type", principal.Id.ResourceType),
			zap.String("principal_id", principal.Id.Resource),
		)
		return nil, fmt.Errorf("baton-verkada: only users can have group membership revoked")
	}

	err := g.client.RemoveUserFromGroup(ctx, entitlement.Resource.Id.Resource, principal.Id.Resource)
	if err != nil {
		return nil, fmt.Errorf("baton-verkada: failed to remove user from group: %w", err)
	}

	return nil, nil
}

func newGroupBuilder(client *verkada.Client) *groupBuilder {
	return &groupBuilder{
		resourceType: groupResourceType,
		client:       client,
	}
}

func groupContainsUser(target string, groups []verkada.Group) bool {
	for _, group := range groups {
		if target == group.GroupID {
			return true
		}
	}
	return false
}
