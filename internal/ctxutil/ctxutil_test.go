package ctxutil

import (
	"context"
	"testing"
)

func TestCtxKey_String(t *testing.T) {
	tests := []struct {
		name string
		ck   CtxKey
		want string
	}{
		{
			name: "CtxKeyTraceID",
			ck:   CtxKeyTraceID,
			want: "X-Trace-Id",
		},
		{
			name: "CtxKeyIdentityID",
			ck:   CtxKeyIdentityID,
			want: "X-Identity-Id",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ck.String(); got != tt.want {
				t.Errorf("CtxKey.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTraceIDFromCtx(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestTraceIDFromCtx 1",
			args: args{
				context.WithValue(context.Background(), CtxKeyTraceID, "abcd-trace-id"),
			},
			want: "abcd-trace-id",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TraceIDFromCtx(tt.args.ctx); got != tt.want {
				t.Errorf("TraceIDFromCtx() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIdentityIDFromCtx(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "TestIdentityIDFromCtx 1",
			args: args{
				context.WithValue(context.Background(), CtxKeyIdentityID, 10),
			},
			want:    10,
			wantErr: false,
		},
		{
			name: "TestIdentityIDFromCtx 2",
			args: args{
				context.WithValue(context.Background(), CtxKeyIdentityID, int64(20)),
			},
			want:    20,
			wantErr: false,
		},
		{
			name: "TestIdentityIDFromCtx 3",
			args: args{
				context.WithValue(context.Background(), CtxKeyIdentityID, "30"),
			},
			want:    30,
			wantErr: false,
		},
		{
			name: "TestIdentityIDFromCtx failed",
			args: args{
				context.Background(),
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IdentityIDFromCtx(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("IdentityIDFromCtx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IdentityIDFromCtx() = %v, want %v", got, tt.want)
			}
		})
	}
}
