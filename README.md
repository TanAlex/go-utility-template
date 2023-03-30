# README

This repo shows how to use go-client to call GCP API to list all firewall rules

All the code was generated using ChatGPT.

## Procedure

```
go mod init github.com/TanAlex/firewall-list-tool
go get github.com/TanAlex/firewall-list-tool

PROJECT_ID="<your project id>"
go run firewall-list.go -projectID $PROJECT_ID
```