curl -v -X PUT "http://localhost:1080/mockserver/expectation" -d '{
  "httpRequest" : {
    "path" : "/v1/node/test_node_id/request/upgrade/test"
  },
  "httpResponse" : {
    "body" : "{\"Name\": \"test\", \"Spec\": \"version: \\\"3\\\"\\nservices:\\n  db:\\n    image: postgres:10.8\\n    expose:\\n      - 5432\\n    environment:\\n      - POSTGRES_USER=db_user_name\\n      - POSTGRES_PASSWORD=P56FJXc\",\"DeleteVolumes\": \"false\",\"HealthCheckCmds\": [{\"ContainerName\": \"db\", \"Cmd\": \"pg_isready -U postgres\"}, {\"ContainerName\": \"cache\", \"Cmd\": \"pg_isready -U postgres\"}]}"
  }
}'




#curl -v -X PUT "http://localhost:1080/mockserver/expectation" -d '{
#  "httpRequest" : {
#    "body" : {
#      "type" : "JSON",
#      "json" : "[{\"ID\": \"container-id\", \"CREATED\": \"DD.MM.YYYY\", \"NAME\": \"boosteroid-web_angular_1\", \"IMAGE\": \"node:9.11\"}, {\"ID\": \"container-id\", \"CREATED\": \"DD.MM.YYYY\", \"NAME\": \"boosteroid-web_django_1\", \"IMAGE\": \"boosteroid-web_django:latest\"}, {\"ID\": \"container-id\", \"CREATED\": \"DD.MM.YYYY\", \"NAME\": \"boosteroid-web_db_1\", \"IMAGE\": \"postgres:10.3\"}]",
#      "matchType" : "STRICT"
#    }
#  },
#  "httpResponse" : {
#    "statusCode" : 200,
#    "body" : "[{\"ID\": \"container-id\", \"NAME\": \"boosteroid-web_db_1\", \"IMAGE\": \"postgres:10.8\", \"DELETE_VOLUMES\": false}]"
#  }
#}'