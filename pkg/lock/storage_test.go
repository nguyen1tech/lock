package lock

import (
	"context"
	"testing"
)

func TestGet(t *testing.T) {
	type testcase struct {
		m         map[string]interface{}
		key       string
		value     interface{}
		wantErr   error
		wantValue interface{}
	}

	var testcases = []testcase{
		{
			m:         map[string]interface{}{},
			key:       "k1",
			value:     "v1",
			wantErr:   ErrKeyNotExisted,
			wantValue: nil,
		}, {
			m:         map[string]interface{}{"k1": "v1", "k2": "v2", "k3": "v3"},
			key:       "k1",
			value:     "v1",
			wantErr:   nil,
			wantValue: "v1",
		},
	}

	for _, tc := range testcases {
		ctx := context.TODO()
		inmem := newInMemoryStorage(tc.m)
		gotErr, gotValue := inmem.Get(ctx, tc.key)
		if gotErr != tc.wantErr {
			t.Errorf("want err: %+v but got: %+v", tc.wantErr, gotErr)
		}
		if gotValue != tc.wantValue {
			t.Errorf("want value: %+v but got: %+v", tc.wantValue, gotValue)
		}
	}
}

func TestSet(t *testing.T) {
	type testcase struct {
		m         map[string]interface{}
		key       string
		value     interface{}
		wantErr   error
		wantValue string
	}

	var testcases = []testcase{
		{
			m:         map[string]interface{}{},
			key:       "k1",
			value:     "v1",
			wantErr:   nil,
			wantValue: "v1",
		}, {
			m:         map[string]interface{}{"k1": "v1"},
			key:       "k1",
			value:     "v2",
			wantErr:   nil,
			wantValue: "v2",
		},
	}

	for _, tc := range testcases {
		ctx := context.TODO()
		inmem := newInMemoryStorage(tc.m)
		gotErr := inmem.Set(ctx, tc.key, tc.value)
		if gotErr != tc.wantErr {
			t.Errorf("want err: %+v but got: %+v", tc.wantErr, gotErr)
		}
		_, gotValue := inmem.Get(ctx, tc.key)
		if gotValue != tc.wantValue {
			t.Errorf("want value: %+v but got: %+v", tc.wantValue, gotValue)
		}
	}
}

func TestDel(t *testing.T) {
	type testcase struct {
		m   map[string]interface{}
		key string
	}

	var testcases = []testcase{
		{
			m:   map[string]interface{}{},
			key: "k1",
		}, {
			m:   map[string]interface{}{"k1": "v1"},
			key: "k1",
		},
	}

	for _, tc := range testcases {
		ctx := context.TODO()
		inmem := newInMemoryStorage(tc.m)
		_ = inmem.Del(ctx, tc.key)
		_, ok := tc.m[tc.key]
		if ok {
			t.Errorf("want key: %+v to be deleted but not", tc.key)
		}
	}
}
