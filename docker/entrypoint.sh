#!/bin/sh

exec /app/sonicd --datadir /app/data \
  --config /app/config.toml \
  --validator.id="$SONIC_VALIDATOR_ID" \
  --validator.pubkey="$SONIC_PUBKEY" \
  --validator.password="/app/key" \
  --nat="$SONIC_NAT" \
  --http \
  --http.addr="0.0.0.0" \
  --http.vhosts="*" \
  --http.api=eth,debug,net,admin,web3,personal,txpool,ftm,dag \
  --bootnodes="$SONIC_BOOT_NODES"
