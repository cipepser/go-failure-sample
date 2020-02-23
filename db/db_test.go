package db

import (
	"github.com/morikuni/failure"
	"reflect"
	"testing"
)

func cleanUpCustomers() {
	customers = nil
}

func TestNewCustomer(t *testing.T) {
	type args struct {
		name    string
		address string
	}
	tests := []struct {
		name string
		args args
		want *Customer
	}{
		{
			name: "single customer",
			args: args{
				name:    "Alice",
				address: "alice@example.com",
			},
			want: &Customer{
				id:      0,
				name:    "Alice",
				address: "alice@example.com",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCustomer(tt.args.name, tt.args.address); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCustomer() = %v, want %v", got, tt.want)
			}
		})
	}
	cleanUpCustomers()
}

func TestShowCustomers(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
		want    []*Customer
	}{
		{
			name:    "Test for Customers in global state. This test should be executed only once",
			wantErr: false,
			want: []*Customer{
				{
					id:      0,
					name:    "Alice",
					address: "alice@example.com",
				},
				{
					id:      1,
					name:    "Bob",
					address: "bob@example.com",
				},
			},
		},
	}
	for _, tt := range tests {
		_ = NewCustomer("Alice", "alice@example.com")
		_ = NewCustomer("Bob", "bob@example.com")

		t.Run(tt.name, func(t *testing.T) {
			if err := ShowCustomers(); (err != nil) != tt.wantErr {
				t.Errorf("ShowCustomers() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
		t.Run(tt.name, func(t *testing.T) {
			if len(customers) != len(tt.want) {
				t.Fatalf("length mismatch. len(customers) = %v, want %v", len(customers), len(tt.want))
			}
			for i := range customers {
				if !reflect.DeepEqual(*customers[i], *tt.want[i]) {
					t.Errorf("customers(global state) = %v, want %v", *customers[i], *tt.want[i])
				}

			}
		})
	}
	cleanUpCustomers()
}

func TestClient_GetName(t *testing.T) {
	type fields struct {
		user string
	}
	type args struct {
		id int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr failure.StringCode
	}{
		{
			name: "Alice",
			args: args{
				id: 0,
			},
			want:    "Alice",
			wantErr: "",
		},
		{
			name: "",
			args: args{
				id: -1,
			},
			want:    "",
			wantErr: NotFound,
		},
	}

	_ = NewCustomer("Alice", "alice@example.com")
	_ = NewCustomer("Bob", "bob@example.com")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				user: tt.fields.user,
			}
			got, err := c.GetName(tt.args.id)
			if err != nil && !failure.Is(err, tt.wantErr) {
				t.Errorf("GetName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetName() got = %v, want %v", got, tt.want)
			}
		})
	}
	cleanUpCustomers()
}
