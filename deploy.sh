name: CI/CD Pipeline

on:
  push:
    branches:
      - main

jobs:
  deploy:
    runs-on: ubuntu-latest

    steps:
    - name: Checkout Code
      uses: actions/checkout@v2

    - name: Set Up Go
      uses: actions/setup-go@v2
      with:
        go-version: 1.16

    - name: Build and Deploy
      run: |
        # Execute your deployment script or command here
        bash deploy.sh
