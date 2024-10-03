package configs

import (
	"os"
	"testing"
)

func TestCompareNetworkConfigs(t *testing.T) {
	type args struct {
		running  string
		designed string
	}
	running, err := os.ReadFile("testdata/running_config.cfg")
	if err != nil {
		t.Fatalf("failed reading file: %s", err)
	}
	designed, err := os.ReadFile("testdata/designed_config.cfg")
	if err != nil {
		t.Fatalf("failed reading file: %s", err)
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test 1",
			args: args{
				running:  string(running),
				designed: string(designed),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CompareNetworkConfigs(tt.args.running, tt.args.designed)
			if got != tt.want {
				t.Errorf("CompareNetworkConfigs() = %v, want %v", got, tt.want)
				t.Logf("Differences found in configurations:")
				CompareNetworkConfigs(tt.args.running, tt.args.designed) // Run again to get the logged output
			}
		})
	}
}
