package token

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/go-logr/logr/testr"
	dockertest "github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"
	"github.com/tomato-net/vault-agent/config"
	"github.com/tomato-net/vault-agent/credentials"
)

func setupTestVault(t *testing.T) string {
	pool, err := dockertest.NewPool("")
	require.NoError(t, err, "could not connect to Docker")

	resource, err := pool.RunWithOptions(
		&dockertest.RunOptions{
			Repository: "vault",
			Tag:        "latest",
			CapAdd:     []string{"IPC_LOCK"},
			Env:        []string{"VAULT_DEV_ROOT_TOKEN_ID=myroot"},
		}, func(hc *docker.HostConfig) {
			// hc.AutoRemove = true
		})
	require.NoError(t, err, "could not start container")

	// t.Cleanup(func() {
	// 	require.NoError(t, pool.Purge(resource), "failed to remove container")
	// })

	var resp *http.Response

	err = pool.Retry(func() error {
		resp, err = http.Get(fmt.Sprint("http://localhost:", resource.GetPort("8200/tcp"), "/ui"))
		if err != nil {
			t.Logf("container not ready: %s, waiting...", err.Error())
			return err
		}
		return nil
	})
	require.NoError(t, err, "HTTP error")
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "HTTP status code")

	exitCode, err := resource.Exec([]string{"vault login myroot", "vault auth enable userpass", "vault write auth/userpass/users/example password=foo policies=admins"}, dockertest.ExecOptions{
		Env: []string{"VAULT_ADDR='http://127.0.0.1:8200'"},
		TTY: true,
	})
	require.Zero(t, exitCode, "vault setup non-zero exit code")
	require.NoError(t, err, "vault setup err")

	return resource.GetPort("8200/tcp")
}

func TestVault(t *testing.T) {
	vaultPort := setupTestVault(t)
	cfg := config.NewFake(fmt.Sprint("http://localhost:", vaultPort), "example", ".my-token")
	client, err := NewClient(cfg)
	require.NoError(t, err, "Client setup error")
	provider := NewProviderUserPass(client, testr.New(t), cfg, &credentials.Fake{Password: "foo"})
	renewer := NewRenewer(client, testr.New(t), provider, cfg)

	go func() {
		gotErr := renewer.Start()
		require.NoError(t, gotErr, "renewer start error")
	}()

	time.Sleep(5 * time.Second)
	require.NoError(t, renewer.Stop(), "renewer stop error")
}
