package node

import (
	"context"
	"crypto/rand"
	"testing"

	"github.com/libp2p/go-libp2p-core/crypto"
	mocknet "github.com/libp2p/go-libp2p/p2p/net/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLight(t *testing.T) {
	store := MockStore(t, DefaultConfig(Light))
	nd, err := New(Light, store)
	require.NoError(t, err)
	require.NotNil(t, nd)
	require.NotNil(t, nd.Config)
	require.NotNil(t, nd.HeaderServ)
	assert.NotZero(t, nd.Type)
}

func TestLightLifecycle(t *testing.T) {
	store := MockStore(t, DefaultConfig(Light))
	nd, err := New(Light, store)
	require.NoError(t, err)

	startCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = nd.Start(startCtx)
	require.NoError(t, err)

	stopCtx, stopCtxCancel := context.WithCancel(context.Background())
	t.Cleanup(func() {
		cancel()
		stopCtxCancel()
	})

	err = nd.Stop(stopCtx)
	require.NoError(t, err)
}

func TestNewLightWithP2PKey(t *testing.T) {
	key, _, err := crypto.GenerateEd25519Key(rand.Reader)
	require.NoError(t, err)

	repo := MockStore(t, DefaultConfig(Light))
	node, err := New(Light, repo, WithP2PKey(key))
	require.NoError(t, err)
	assert.True(t, node.Host.ID().MatchesPrivateKey(key))
}

func TestNewLightWithHost(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	nw, _ := mocknet.WithNPeers(ctx, 1)
	repo := MockStore(t, DefaultConfig(Light))
	node, err := New(Light, repo, WithHost(nw.Host(nw.Peers()[0])))
	require.NoError(t, err)
	assert.Equal(t, node.Host.ID(), nw.Peers()[0])
}
