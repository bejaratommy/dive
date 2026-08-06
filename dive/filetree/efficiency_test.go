package filetree

import (
	"testing"
)

func checkError(t *testing.T, err error, message string) {
	if err != nil {
		t.Errorf(message+": %+v", err)
	}
}

func TestEfficiency(t *testing.T) {
	trees := make([]*FileTree, 3)
	for idx := range trees {
		trees[idx] = NewFileTree()
	}

	_, _, err := trees[0].AddPath("/etc/nginx/nginx.conf", FileInfo{Size: 2000})
	checkError(t, err, "could not setup test")

	_, _, err = trees[0].AddPath("/etc/nginx/public", FileInfo{Size: 3000})
	checkError(t, err, "could not setup test")

	_, _, err = trees[1].AddPath("/etc/nginx/nginx.conf", FileInfo{Size: 5000})
	checkError(t, err, "could not setup test")
	_, _, err = trees[1].AddPath("/etc/athing", FileInfo{Size: 10000})
	checkError(t, err, "could not setup test")

	_, _, err = trees[2].AddPath("/etc/.wh.nginx", *BlankFileChangeInfo("/etc/.wh.nginx"))
	checkError(t, err, "could not setup test")

	var expectedScore = 0.75
	var expectedMatches = EfficiencySlice{
		&EfficiencyData{Path: "/etc/nginx/nginx.conf", CumulativeSize: 7000, WastedSize: 5000},
	}
	actualScore, actualMatches := Efficiency(trees)

	if expectedScore != actualScore {
		t.Errorf("Expected score of %v but go %v", expectedScore, actualScore)
	}

	if len(actualMatches) != len(expectedMatches) {
		for _, match := range actualMatches {
			t.Logf("   match: %+v", match)
		}
		t.Fatalf("Expected to find %d inefficient paths, but found %d", len(expectedMatches), len(actualMatches))
	}

	if expectedMatches[0].Path != actualMatches[0].Path {
		t.Errorf("Expected path of %s but go %s", expectedMatches[0].Path, actualMatches[0].Path)
	}

	if expectedMatches[0].CumulativeSize != actualMatches[0].CumulativeSize {
		t.Errorf("Expected cumulative size of %v but go %v", expectedMatches[0].CumulativeSize, actualMatches[0].CumulativeSize)
	}

	if expectedMatches[0].WastedSize != actualMatches[0].WastedSize {
		t.Errorf("Expected wasted size of %v but go %v", expectedMatches[0].WastedSize, actualMatches[0].WastedSize)
	}
}

// TestEfficiency_WastedSize covers the reclaimable bytes reported for a path,
// which is every copy of that path but one. A path that only ever appears once
// wastes nothing, even though it does contribute to the cumulative size.
func TestEfficiency_WastedSize(t *testing.T) {
	tests := []struct {
		name               string
		layers             [][]string
		size               int64
		path               string
		expectedWastedSize int64
	}{
		{
			name:               "single copy wastes nothing",
			layers:             [][]string{{"/usr/bin/app"}},
			size:               1000,
			path:               "/usr/bin/app",
			expectedWastedSize: 0,
		},
		{
			name:               "two copies waste one",
			layers:             [][]string{{"/usr/bin/app"}, {"/usr/bin/app"}},
			size:               1000,
			path:               "/usr/bin/app",
			expectedWastedSize: 1000,
		},
		{
			name:               "four copies waste three",
			layers:             [][]string{{"/usr/bin/app"}, {"/usr/bin/app"}, {"/usr/bin/app"}, {"/usr/bin/app"}},
			size:               1000,
			path:               "/usr/bin/app",
			expectedWastedSize: 3000,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			trees := make([]*FileTree, len(test.layers))
			for idx, paths := range test.layers {
				trees[idx] = NewFileTree()
				for _, path := range paths {
					_, _, err := trees[idx].AddPath(path, FileInfo{Size: test.size})
					checkError(t, err, "could not setup test")
				}
			}

			_, matches := Efficiency(trees)

			var actualWastedSize int64
			for _, match := range matches {
				if match.Path == test.path {
					actualWastedSize = match.WastedSize
				}
			}

			if actualWastedSize != test.expectedWastedSize {
				t.Errorf("Expected wasted size of %v but go %v", test.expectedWastedSize, actualWastedSize)
			}
		})
	}
}

func TestEfficiency_ScratchImage(t *testing.T) {
	trees := make([]*FileTree, 3)
	for idx := range trees {
		trees[idx] = NewFileTree()
	}

	_, _, err := trees[0].AddPath("/nothing", FileInfo{Size: 0})
	checkError(t, err, "could not setup test")

	var expectedScore = 1.0
	var expectedMatches = EfficiencySlice{}
	actualScore, actualMatches := Efficiency(trees)

	if expectedScore != actualScore {
		t.Errorf("Expected score of %v but go %v", expectedScore, actualScore)
	}

	if len(actualMatches) > 0 {
		t.Fatalf("Expected to find %d inefficient paths, but found %d", len(expectedMatches), len(actualMatches))
	}

}
