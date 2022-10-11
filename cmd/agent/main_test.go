package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_genURL(t *testing.T) {
	type args struct {
		host        string
		port        string
		metricType  string
		metricName  string
		metricValue string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty args",
			args: args{
				host:        "",
				port:        "",
				metricType:  "",
				metricName:  "",
				metricValue: "",
			},
			want: "http://:/update///",
		},
		{
			name: "correct args",
			args: args{
				host:        host,
				port:        port,
				metricValue: "1.321456",
				metricName:  "Alloc",
				metricType:  "gauge",
			},
			want: fmt.Sprintf("http://%v:%v/update/gauge/Alloc/1.321456", host, port),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := genURL(tt.args.host, tt.args.port, tt.args.metricType, tt.args.metricName, tt.args.metricValue)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_genRandomGauge(t *testing.T) {
	tests := []struct {
		name string
		want gauge
	}{
		{
			name: "different results",
			want: genRandomGauge(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := genRandomGauge()
			assert.NotEqual(t, tt.want, got)
		})
	}
}
