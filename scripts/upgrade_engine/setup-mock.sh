curl -v -X PUT "http://localhost:1080/mockserver/expectation" -d '{
  "httpRequest" : {
    "body" : {
      "type" : "JSON",
      "json" : "[{\"ID\": \"container-id\", \"CREATED\": \"DD.MM.YYYY\", \"NAME\": \"boosteroid-web_angular_1\", \"IMAGE\": \"node:9.11\"}, {\"ID\": \"container-id\", \"CREATED\": \"DD.MM.YYYY\", \"NAME\": \"boosteroid-web_django_1\", \"IMAGE\": \"boosteroid-web_django:latest\"}, {\"ID\": \"container-id\", \"CREATED\": \"DD.MM.YYYY\", \"NAME\": \"boosteroid-web_db_1\", \"IMAGE\": \"postgres:10.3\"}]",
      "matchType" : "STRICT"
    }
  },
  "httpResponse" : {
    "statusCode" : 200,
    "body" : "[{\"ID\": \"container-id\", \"NAME\": \"boosteroid-web_db_1\", \"IMAGE\": \"postgres:10.8\", \"DELETE_VOLUMES\": false}]"
  }
}'