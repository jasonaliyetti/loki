package loki

import (
	"testing"
	"time"

	"github.com/cortexproject/cortex/pkg/chunk"
)

func Test_calculateMaxLookBack(t *testing.T) {
	type args struct {
		pc                               chunk.PeriodConfig
		maxLookBackConfig                time.Duration
		maxChunkAge                      time.Duration
		querierBoltDBFilesResyncInterval time.Duration
	}
	tests := []struct {
		name    string
		args    args
		want    time.Duration
		wantErr bool
	}{
		{
			name: "default",
			args: args{
				pc: chunk.PeriodConfig{
					ObjectType: "filesystem",
				},
				maxLookBackConfig:                0,
				maxChunkAge:                      1 * time.Hour,
				querierBoltDBFilesResyncInterval: 5 * time.Minute,
			},
			want:    81 * time.Minute,
			wantErr: false,
		},
		{
			name: "infinite",
			args: args{
				pc: chunk.PeriodConfig{
					ObjectType: "filesystem",
				},
				maxLookBackConfig:                -1,
				maxChunkAge:                      1 * time.Hour,
				querierBoltDBFilesResyncInterval: 5 * time.Minute,
			},
			want:    -1,
			wantErr: false,
		},
		{
			name: "invalid store type",
			args: args{
				pc: chunk.PeriodConfig{
					ObjectType: "gcs",
				},
				maxLookBackConfig:                -1,
				maxChunkAge:                      1 * time.Hour,
				querierBoltDBFilesResyncInterval: 5 * time.Minute,
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "less than default",
			args: args{
				pc: chunk.PeriodConfig{
					ObjectType: "filesystem",
				},
				maxLookBackConfig:                1 * time.Hour,
				maxChunkAge:                      1 * time.Hour,
				querierBoltDBFilesResyncInterval: 5 * time.Minute,
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculateMaxLookBack(tt.args.pc, tt.args.maxLookBackConfig, tt.args.maxChunkAge, tt.args.querierBoltDBFilesResyncInterval)
			if (err != nil) != tt.wantErr {
				t.Errorf("calculateMaxLookBack() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("calculateMaxLookBack() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_calculateAsyncStoreQueryIngestersWithin(t *testing.T) {
	type args struct {
		queryIngestersWithin             time.Duration
		maxChunkAge                      time.Duration
		querierBoltDBFilesResyncInterval time.Duration
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{
			name: "default",
			args: args{
				0,
				time.Hour,
				5 * time.Minute,
			},
			want: 0,
		},
		{
			name: "queryIngestersWithin more than min val",
			args: args{
				3 * time.Hour,
				time.Hour,
				5 * time.Minute,
			},
			want: 3 * time.Hour,
		},
		{
			name: "queryIngestersWithin less than min val",
			args: args{
				time.Hour,
				time.Hour,
				5 * time.Minute,
			},
			want: 81 * time.Minute,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateAsyncStoreQueryIngestersWithin(tt.args.queryIngestersWithin, tt.args.maxChunkAge, tt.args.querierBoltDBFilesResyncInterval)

			if got != tt.want {
				t.Errorf("calculateAsyncStoreQueryIngestersWithin() got = %v, want %v", got, tt.want)
			}
		})
	}
}
