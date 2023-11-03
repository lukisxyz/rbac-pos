#!/bin/bash

# Define the server URL and Bearer token
SERVER_URL="http://localhost:8080"
BEARER_TOKEN=$1

# Function to create a role
create_role() {
  name="$1"
  description="$2"
  data='{"name":"'$name'", "description":"'$description'"}'
  curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $BEARER_TOKEN" -d "$data" "$SERVER_URL/api/role"
}

# Function to create a permission
create_permission() {
  name="$1"
  description="$2"
  url="$3"
  data='{"name":"'$name'", "description":"'$description'", "url":"'$url'"}'
  curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $BEARER_TOKEN" -d "$data" "$SERVER_URL/api/permission"
}

# Create 3 roles
create_role "Cashier" "Responsible for processing sales and refunds."
create_role "Manager" "Has additional access to inventory management and reports."
create_role "Administrator" "Full access to all system features and settings."

# Create 20 permissions with real-world examples
create_permission "Create Sale" "Permission to create a new sale transaction." "create-sale"
create_permission "Edit Sale" "Permission to modify an existing sale transaction." "edit-sale"
create_permission "Refund Transaction" "Permission to process refunds for sales." "refund-transaction"
create_permission "View Inventory" "Permission to view inventory and stock levels." "view-inventory"
create_permission "Manage Inventory" "Permission to add, update, or remove items from inventory." "manage-inventory"
create_permission "Generate Reports" "Permission to generate sales and inventory reports." "generate-reports"
create_permission "Customer Management" "Permission to add, edit, or delete customer records." "customer-management"
create_permission "User Management" "Permission to manage user accounts and roles." "user-management"
create_permission "Access Settings" "Permission to access and configure application settings." "access-settings"
