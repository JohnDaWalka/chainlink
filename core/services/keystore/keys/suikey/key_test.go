package suikey

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
)

func Test_SuiKeystore(t *testing.T) {
	ctx := context.Background()
	ks := keystore.NewInMemory()
	suiKs := New(ks)

	t.Run("Create", func(t *testing.T) {
		account, err := suiKs.Create(ctx)
		require.NoError(t, err)
		assert.NotNil(t, account)
		assert.NotEmpty(t, account.ID())
	})

	t.Run("Get", func(t *testing.T) {
		account, err := suiKs.Create(ctx)
		require.NoError(t, err)

		got, err := suiKs.Get(account.ID())
		require.NoError(t, err)
		assert.Equal(t, account.ID(), got.ID())
	})

	t.Run("GetAll", func(t *testing.T) {
		_, err := suiKs.Create(ctx)
		require.NoError(t, err)
		_, err = suiKs.Create(ctx)
		require.NoError(t, err)

		accounts, err := suiKs.GetAll()
		require.NoError(t, err)
		assert.Len(t, accounts, 2)
	})

	t.Run("Delete", func(t *testing.T) {
		account, err := suiKs.Create(ctx)
		require.NoError(t, err)

		err = suiKs.Delete(account.ID())
		require.NoError(t, err)

		_, err = suiKs.Get(account.ID())
		assert.Error(t, err)
	})

	t.Run("Import/Export", func(t *testing.T) {
		account, err := suiKs.Create(ctx)
		require.NoError(t, err)

		exported, err := suiKs.Export(account.ID(), "password")
		require.NoError(t, err)

		imported, err := suiKs.Import(exported, "password")
		require.NoError(t, err)
		assert.Equal(t, account.ID(), imported.ID())
	})
}
