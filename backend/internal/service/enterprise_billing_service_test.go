//go:build unit

package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Tests — EnterpriseBillingService
// ---------------------------------------------------------------------------

func TestBilling_GetFinanceOverview(t *testing.T) {
	ent := &entRepoStubForMember{
		enterprise: &Enterprise{ID: 1, Name: "Corp", Balance: 5000, TotalRecharged: 10000, Status: EnterpriseStatusActive},
	}
	sub := &subRepoStubForBilling{
		subs: []EnterpriseSubscription{
			{ID: 1, EnterpriseID: 1, Status: EnterpriseSubStatusActive},
		},
	}
	svc := NewEnterpriseBillingService(ent, sub, &memberRepoStubForService{})

	overview, err := svc.GetFinanceOverview(context.Background(), 1)
	require.NoError(t, err)
	require.Equal(t, float64(5000), overview.Balance)
	require.Equal(t, float64(10000), overview.TotalRecharged)
	require.Len(t, overview.Subscriptions, 1)
	require.NotNil(t, overview.MonthlyUsage)
}

func TestBilling_GetSubscriptions(t *testing.T) {
	sub := &subRepoStubForBilling{
		subs: []EnterpriseSubscription{{ID: 1}, {ID: 2}},
	}
	svc := NewEnterpriseBillingService(&entRepoStubForMember{}, sub, &memberRepoStubForService{})

	result, err := svc.GetSubscriptions(context.Background(), 1)
	require.NoError(t, err)
	require.Len(t, result, 2)
}

func TestBilling_GetActiveSubscriptions(t *testing.T) {
	sub := &subRepoStubForBilling{
		subs: []EnterpriseSubscription{{ID: 1, Status: EnterpriseSubStatusActive}},
	}
	svc := NewEnterpriseBillingService(&entRepoStubForMember{}, sub, &memberRepoStubForService{})

	result, err := svc.GetActiveSubscriptions(context.Background(), 1)
	require.NoError(t, err)
	require.Len(t, result, 1)
}

func TestBilling_ValidateRecharge_Success(t *testing.T) {
	svc := NewEnterpriseBillingService(&entRepoStubForMember{}, &subRepoStubForBilling{}, &memberRepoStubForService{})
	err := svc.ValidateRecharge(EnterpriseRechargeInput{Amount: 100})
	require.NoError(t, err)
}

func TestBilling_ValidateRecharge_ZeroAmount(t *testing.T) {
	svc := NewEnterpriseBillingService(&entRepoStubForMember{}, &subRepoStubForBilling{}, &memberRepoStubForService{})
	err := svc.ValidateRecharge(EnterpriseRechargeInput{Amount: 0})
	require.ErrorIs(t, err, ErrEnterpriseRechargeAmountInvalid)
}

func TestBilling_ValidateRecharge_NegativeAmount(t *testing.T) {
	svc := NewEnterpriseBillingService(&entRepoStubForMember{}, &subRepoStubForBilling{}, &memberRepoStubForService{})
	err := svc.ValidateRecharge(EnterpriseRechargeInput{Amount: -10})
	require.ErrorIs(t, err, ErrEnterpriseRechargeAmountInvalid)
}

func TestBilling_GetBalance(t *testing.T) {
	entWithBal := &entRepoStub{balance: 9999.99}
	sub := &subRepoStub{}
	svc := NewEnterpriseBillingService(entWithBal, sub, &memberRepoStub{})
	bal, err := svc.GetBalance(context.Background(), 1)
	require.NoError(t, err)
	require.Equal(t, 9999.99, bal)
}

// ---------------------------------------------------------------------------
// Stub for billing — subRepoStub reuses ent_service_test but needs GetByEnterprise
// ---------------------------------------------------------------------------

type subRepoStubForBilling struct {
	subs          []EnterpriseSubscription
	listActiveErr error
	getByEntErr   error
}

func (s *subRepoStubForBilling) Create(ctx context.Context, sub *EnterpriseSubscription) error { panic("unexpected") }
func (s *subRepoStubForBilling) GetByID(ctx context.Context, id int64) (*EnterpriseSubscription, error) { panic("unexpected") }
func (s *subRepoStubForBilling) GetByEnterprise(ctx context.Context, enterpriseID int64) ([]EnterpriseSubscription, error) {
	if s.getByEntErr != nil {
		return nil, s.getByEntErr
	}
	return s.subs, nil
}
func (s *subRepoStubForBilling) ListActive(ctx context.Context, enterpriseID int64) ([]EnterpriseSubscription, error) {
	if s.listActiveErr != nil {
		return nil, s.listActiveErr
	}
	return s.subs, nil
}
func (s *subRepoStubForBilling) UpdateStatus(ctx context.Context, id int64, status string) error { panic("unexpected") }
