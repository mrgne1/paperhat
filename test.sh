resp1=$(curl -d "my secret value" -X POST localhost:2060/api/secrets)
echo $resp1
resp2=$(curl -X GET localhost:2060/api/secrets/$resp1)
echo $resp2
resp3=$(curl -X GET localhost:2060/api/secrets/$resp1)
echo $resp3
