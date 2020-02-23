package db

import (
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
