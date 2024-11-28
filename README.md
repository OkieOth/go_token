# go_token
Golang library to handle keycloak tokens


# Curl helper

```shell
curl --request POST \
  --url http://localhost:8080/realms/test-realm/protocol/openid-connect/token \
  --header 'Content-Type: application/x-www-form-urlencoded' \
  --data 'client_id=test-client' \
  --data 'grant_type=client_credentials' \
  --data 'client_secret=test-client999'

```

