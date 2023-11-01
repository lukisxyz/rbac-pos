#!/bin/bash

# Define the server URL
SERVER_URL="http://localhost:8080"

# Function to create an account
create_account() {
  email="$1"
  password="$2"
  data='{"email":"'$email'", "password":"'$password'"}'
  curl -X POST -H "Content-Type: application/json" -d "$data" "$SERVER_URL/api/account"
}

# Function to create a role
create_role() {
  name="$1"
  description="$2"
  data='{"name":"'$name'", "description":"'$description'"}'
  curl -X POST -H "Content-Type: application/json" -d "$data" "$SERVER_URL/api/role"
}

# Function to create a permission
create_permission() {
  name="$1"
  description="$2"
  url="$3"
  data='{"name":"'$name'", "description":"'$description'", "url":"'$url'"}'
  curl -X POST -H "Content-Type: application/json" -d "$data" "$SERVER_URL/api/permission"
}

# Create 10 accounts
for i in {1..10}; do
  create_account "user$i@example.com" "password$i"
done

# Create 3 roles
create_role "Cashier" "Responsible for processing sales and refunds."
create_role "Manager" "Has additional access to inventory management and reports."
create_role "Administrator" "Full access to all system features and settings."

# Create 20 permissions with real-world examples
create_permission "Create Sale" "Permission to create a new sale transaction." "/api/create-sale"
create_permission "Edit Sale" "Permission to modify an existing sale transaction." "/api/edit-sale"
create_permission "Refund Transaction" "Permission to process refunds for sales." "/api/refund-transaction"
create_permission "View Inventory" "Permission to view inventory and stock levels." "/api/view-inventory"
create_permission "Manage Inventory" "Permission to add, update, or remove items from inventory." "/api/manage-inventory"
create_permission "Generate Reports" "Permission to generate sales and inventory reports." "/api/generate-reports"
create_permission "Customer Management" "Permission to add, edit, or delete customer records." "/api/customer-management"
create_permission "User Management" "Permission to manage user accounts and roles." "/api/user-management"
create_permission "Access Settings" "Permission to access and configure application settings." "/api/access-settings"
