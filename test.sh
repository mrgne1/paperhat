secret="my secret value"
resp1=$(curl -d "$secret" -X POST localhost:2060/api/secrets)
echo $resp1
echo $(echo $resp1 | jq -r .directLink)
resp2=$(curl -X GET $(echo $resp1 | jq -r .directLink ))
echo $resp2
resp3=$(curl -X GET $(echo $resp1 | jq -r .directLink))
echo $resp3
