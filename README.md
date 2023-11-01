# RBAC POS Using Golang
Golang RBAC (role based access controll) on POS application

## Introduction

- Account 1 - 1 Role
- Role 1 - Permission
- Permission 1 - 1 Protected

We need atleas some of this table:
- Account table
- Session table
- Role table
- Permission table
- Role - permission table
- Account - role table

- Protected (optionals) <-  no need, just return authorized or unauthorized
