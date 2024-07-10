package glicko2

import (
	"reflect"
	"testing"

	"matcher/common"
	"matcher/pto"
)

func TestNewPlayer(t *testing.T) {
	type args struct {
		base *common.Base
	}
	tests := []struct {
		name    string
		args    args
		want    *Player
		wantErr bool
	}{
		{
			name: "base is nil should failed",
			args: args{
				base: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "base is not nil should work",
			args: args{
				base: common.NewBase(&pto.PlayerInfo{}),
			},
			want:    &Player{Base: common.NewBase(&pto.PlayerInfo{})},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPlayer(tt.args.base)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPlayer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPlayer() got = %v, want %v", got, tt.want)
			}
		})
	}
}
