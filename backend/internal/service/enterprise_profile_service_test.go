//go:build unit

package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Tests — EnterpriseProfileService
// ---------------------------------------------------------------------------

func TestProfile_GetProfile_Success(t *testing.T) {
	ent := &entRepoStubForMember{
		enterprise: &Enterprise{ID: 1, Name: "My Corp", Status: EnterpriseStatusActive},
	}
	mem := &memberRepoStubForService{
		member: &EnterpriseMember{ID: 1, EnterpriseID: 1, UserID: 500, Role: EnterpriseRoleMember, Status: StatusActive},
	}
	dept := &deptRepoStubForMember{dept: &Department{ID: 1, Name: "Eng", EnterpriseID: 1}}
	svc := NewEnterpriseProfileService(ent, mem, dept)

	profile, err := svc.GetProfile(context.Background(), 500)
	require.NoError(t, err)
	require.Equal(t, "My Corp", profile.Enterprise.Name)
	require.NotNil(t, profile.Member)
	require.NotNil(t, profile.MonthlyUsage)
	// Monthly usage is 0 during P3 (TODO P4)
	require.Equal(t, int64(0), profile.MonthlyUsage.CallCount)
}

func TestProfile_GetProfile_NotMember(t *testing.T) {
	mem := &memberRepoStubForService{
		member: &EnterpriseMember{ID: 1, EnterpriseID: 1, UserID: 500, Status: StatusDisabled},
	}
	svc := NewEnterpriseProfileService(&entRepoStubForMember{}, mem, &deptRepoStubForMember{})

	_, err := svc.GetProfile(context.Background(), 500)
	require.ErrorIs(t, err, ErrProfileNotEnterpriseMember)
}

func TestProfile_GetProfileByEnterprise_Success(t *testing.T) {
	ent := &entRepoStubForMember{
		enterprise: &Enterprise{ID: 1, Name: "Corp A", Status: EnterpriseStatusActive},
	}
	svc := NewEnterpriseProfileService(ent, &memberRepoStubForService{}, &deptRepoStubForMember{})

	profile, err := svc.GetProfileByEnterprise(context.Background(), 1)
	require.NoError(t, err)
	require.Equal(t, "Corp A", profile.Enterprise.Name)
	require.Nil(t, profile.Member) // admin view, no member context
}
