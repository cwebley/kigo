#!/bin/sh

## 

echo "Winning Player:"
read WINNER
echo "Winning Character:"
read WXTER
echo "Losing Player:"
read LOSER
echo "Losing Character:"
read LXTER

DATA='{"results":{"winner":{"player":"'$WINNER'","xter":"'$WXTER'"},"loser":{"player":"'$LOSER'","xter":"'$LXTER'"}}}'

echo $DATA 

curl -XPOST -d '{"winner":{"player":"'$WINNER'","xter":"'$WXTER'"},"loser":{"player":"'$LOSER'","xter":"'$LXTER'"}}' -H "Content-Type:application/json" http://localhost:8999/update
