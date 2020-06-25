package handler_test

import (
	"testing"

	"github.com/ovrclk/akash/testutil"
	"github.com/ovrclk/akash/x/deployment/handler"
	"github.com/ovrclk/akash/x/deployment/types"
	mtypes "github.com/ovrclk/akash/x/market/types"
	"github.com/stretchr/testify/assert"
)

func TestEndBlock(t *testing.T) {
	suite := setupTestSuite(t)
	d0 := testutil.Deployment(suite.t)
	g0 := testutil.DeploymentGroup(suite.t, d0.DeploymentID, uint32(5))
	g1 := testutil.DeploymentGroup(suite.t, d0.DeploymentID, uint32(100))
	//g1.State = types.GroupClosed

	d1 := testutil.Deployment(suite.t)
	d1.State = types.DeploymentClosed
	g2 := testutil.DeploymentGroup(suite.t, d1.DeploymentID, uint32(8))

	// create deployments in storage
	df := func(s *testSuite, d types.Deployment, groups ...types.Group) {
		grps := make([]types.GroupSpec, 0)
		for _, g := range groups {
			t.Logf("%#v", g.GroupSpec)
			grps = append(grps, g.GroupSpec)
		}
		m := types.MsgCreateDeployment{
			ID:     d.ID(),
			Groups: grps,
		}
		_, err := suite.handler(suite.ctx, m)
		assert.NoError(suite.t, err)
	}
	df(suite, d0, g0, g1)
	df(suite, d1, g2)

	handler.OnEndBlock(suite.ctx, suite.dkeeper, suite.mkeeper)

	gx := suite.dkeeper.GetGroups(suite.ctx, d0.ID())
	if len(gx) == 0 {
		t.Error("no groups returned from keeper")
	}
	for _, g := range gx {
		orderCreated := false
		suite.mkeeper.WithOrdersForGroup(suite.ctx, g.ID(), func(o mtypes.Order) bool {
			suite.t.Logf("Order for group: %#v found", o.GroupID())
			orderCreated = true
			return true
		})
		if !orderCreated {
			suite.t.Error("order was not created for a group")
		}
	}
}
