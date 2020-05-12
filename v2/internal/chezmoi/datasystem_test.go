package chezmoi

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twpayne/go-vfs/vfst"
)

var _ System = &DataSystem{}

func TestDataSystem(t *testing.T) {
	fs, cleanup, err := vfst.NewTestFS(map[string]interface{}{
		"/home/user/.local/share/chezmoi": map[string]interface{}{
			".chezmoiignore":  "README.md\n",
			".chezmoiremove":  "*.txt\n",
			".chezmoiversion": "1.2.3\n",
			".chezmoitemplates": map[string]interface{}{
				"foo": "bar",
			},
			"README.md": "",
			"dir": map[string]interface{}{
				"foo": "bar",
			},
			"run_script":      "#!/bin/sh\n",
			"symlink_symlink": "bar",
		},
	})
	require.NoError(t, err)
	defer cleanup()

	s := NewSourceState(
		WithSystem(NewRealSystem(fs, newTestPersistentState())),
		WithSourcePath("/home/user/.local/share/chezmoi"),
	)
	require.NoError(t, s.Read())
	require.NoError(t, s.Evaluate())

	dataSystem := NewDataSystem()
	require.NoError(t, s.ApplyAll(dataSystem, vfst.DefaultUmask, ""))
	expectedData := map[string]interface{}{
		"dir": &dirData{
			Type: dataTypeDir,
			Name: "dir",
			Perm: 0o755,
		},
		"dir/foo": &fileData{
			Type:     dataTypeFile,
			Name:     "dir/foo",
			Contents: "bar",
			Perm:     0o644,
		},
		"script": &scriptData{
			Type:     dataTypeScript,
			Name:     "script",
			Contents: "#!/bin/sh\n",
		},
		"symlink": &symlinkData{
			Type:     dataTypeSymlink,
			Name:     "symlink",
			Linkname: "bar",
		},
	}
	actualData := dataSystem.Data()
	assert.Equal(t, expectedData, actualData)
}