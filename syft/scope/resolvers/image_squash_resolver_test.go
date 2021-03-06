package resolvers

import (
	"testing"

	"github.com/anchore/stereoscope/pkg/imagetest"

	"github.com/anchore/stereoscope/pkg/file"
)

func TestImageSquashResolver_FilesByPath(t *testing.T) {
	cases := []struct {
		name         string
		linkPath     string
		resolveLayer uint
		resolvePath  string
	}{
		{
			name:         "link with previous data",
			linkPath:     "/link-1",
			resolveLayer: 1,
			resolvePath:  "/file-1.txt",
		},
		{
			name:         "link with in layer data",
			linkPath:     "/link-within",
			resolveLayer: 5,
			resolvePath:  "/file-3.txt",
		},
		{
			name:         "link with overridden data",
			linkPath:     "/link-2",
			resolveLayer: 7,
			resolvePath:  "/file-2.txt",
		},
		{
			name:         "indirect link (with overridden data)",
			linkPath:     "/link-indirect",
			resolveLayer: 7,
			resolvePath:  "/file-2.txt",
		},
		{
			name:         "dead link",
			linkPath:     "/link-dead",
			resolveLayer: 8,
			resolvePath:  "/link-dead",
		},
		{
			name:        "ignore directories",
			linkPath:    "/bin",
			resolvePath: "",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			img, cleanup := imagetest.GetFixtureImage(t, "docker-archive", "image-symlinks")
			defer cleanup()

			resolver, err := NewImageSquashResolver(img)
			if err != nil {
				t.Fatalf("could not create resolver: %+v", err)
			}

			refs, err := resolver.FilesByPath(file.Path(c.linkPath))
			if err != nil {
				t.Fatalf("could not use resolver: %+v", err)
			}

			expectedRefs := 1
			if c.resolvePath == "" {
				expectedRefs = 0
			}

			if len(refs) != expectedRefs {
				t.Fatalf("unexpected number of resolutions: %d", len(refs))
			}

			if expectedRefs == 0 {
				// nothing else to assert
				return
			}

			actual := refs[0]

			if actual.Path != file.Path(c.resolvePath) {
				t.Errorf("bad resolve path: '%s'!='%s'", actual.Path, c.resolvePath)
			}

			entry, err := img.FileCatalog.Get(actual)
			if err != nil {
				t.Fatalf("failed to get metadata: %+v", err)
			}

			if entry.Source.Metadata.Index != c.resolveLayer {
				t.Errorf("bad resolve layer: '%d'!='%d'", entry.Source.Metadata.Index, c.resolveLayer)
			}
		})
	}
}

func TestImageSquashResolver_FilesByGlob(t *testing.T) {
	cases := []struct {
		name         string
		glob         string
		resolveLayer uint
		resolvePath  string
	}{
		{
			name:         "link with previous data",
			glob:         "**link-1",
			resolveLayer: 1,
			resolvePath:  "/file-1.txt",
		},
		{
			name:         "link with in layer data",
			glob:         "**link-within",
			resolveLayer: 5,
			resolvePath:  "/file-3.txt",
		},
		{
			name:         "link with overridden data",
			glob:         "**link-2",
			resolveLayer: 7,
			resolvePath:  "/file-2.txt",
		},
		{
			name:         "indirect link (with overridden data)",
			glob:         "**link-indirect",
			resolveLayer: 7,
			resolvePath:  "/file-2.txt",
		},
		{
			name:         "dead link",
			glob:         "**link-dead",
			resolveLayer: 8,
			resolvePath:  "/link-dead",
		},
		{
			name:        "ignore directories",
			glob:        "**/bin",
			resolvePath: "",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			img, cleanup := imagetest.GetFixtureImage(t, "docker-archive", "image-symlinks")
			defer cleanup()

			resolver, err := NewImageSquashResolver(img)
			if err != nil {
				t.Fatalf("could not create resolver: %+v", err)
			}

			refs, err := resolver.FilesByGlob(c.glob)
			if err != nil {
				t.Fatalf("could not use resolver: %+v", err)
			}

			expectedRefs := 1
			if c.resolvePath == "" {
				expectedRefs = 0
			}

			if len(refs) != expectedRefs {
				t.Fatalf("unexpected number of resolutions: %d", len(refs))
			}

			if expectedRefs == 0 {
				// nothing else to assert
				return
			}

			actual := refs[0]

			if actual.Path != file.Path(c.resolvePath) {
				t.Errorf("bad resolve path: '%s'!='%s'", actual.Path, c.resolvePath)
			}

			entry, err := img.FileCatalog.Get(actual)
			if err != nil {
				t.Fatalf("failed to get metadata: %+v", err)
			}

			if entry.Source.Metadata.Index != c.resolveLayer {
				t.Errorf("bad resolve layer: '%d'!='%d'", entry.Source.Metadata.Index, c.resolveLayer)
			}
		})
	}
}
