package httpapi

import (
	"reflect"
	"testing"
)

func TestBuildVersionResponse(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		version string
		commit  string
		want    versionResponse
	}{
		{
			name:    "exact tag",
			version: "v1.2.3",
			commit:  "abc",
			want: versionResponse{
				Version: "v1.2.3", Tag: "v1.2.3", Number: "1.2.3", CommitsAfterTag: 0, Commit: "abc", Dirty: false,
			},
		},
		{
			name:    "describe after tag",
			version: "v1.0.0-4-gdeadbeef",
			commit:  "",
			want: versionResponse{
				Version: "v1.0.0-4-gdeadbeef", Tag: "v1.0.0", Number: "1.0.0", CommitsAfterTag: 4, Commit: "", Dirty: false,
			},
		},
		{
			name:    "dirty",
			version: "v2.0.0-1-gabc-dirty",
			commit:  "short",
			want: versionResponse{
				Version: "v2.0.0-1-gabc-dirty", Tag: "v2.0.0", Number: "2.0.0", CommitsAfterTag: 1, Commit: "short", Dirty: true,
			},
		},
		{
			name:    "dev",
			version: "dev",
			commit:  "",
			want: versionResponse{
				Version: "dev", Tag: "dev", Number: "", CommitsAfterTag: 0, Commit: "", Dirty: false,
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := buildVersionResponse(tc.version, tc.commit)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("want %+v, got %+v", tc.want, got)
			}
		})
	}
}
