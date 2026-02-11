
rewardchaind tx rewardchain create-partner "Acme Inc" "retail" "Mumbai" "IN" \
  --from alice \
  --keyring-backend file \
  --chain-id rewardchain \
  --home ~/.rewardchain \
  --fees 1000ureward \
  --yes

  rewardchaind tx rewardchain update-partner <PARTNER_ID> "Acme Inc" "retail" "Mumbai" "IN" \
  --from alice \
  --keyring-backend file \
  --chain-id rewardchain \
  --home ~/.rewardchain \
  --fees 1000ureward \
  --yes

  rewardchaind tx rewardchain disable-partner <PARTNER_ID> \
  --from alice\
  --keyring-backend file \
  --chain-id rewardchain \
  --home ~/.rewardchain \
  --fees 1000ureward \
  --yes

  rewardchaind query rewardchain partner 
  
  rewardchaind tx rewardchain create-partner \
  "Acme Corporation" \
  "retail" \
  "US" \
  "USD" \
  "0.10" \
  "0.15" \
  "1000000" \
  --from alice \
  --keyring-backend file \
  --chain-id rewardchain \
  --home ~/.rewardchain \
  --fees 1000token \
  --yes


rewardchaind keys add alice \
  --keyring-backend file \
  --home ~/.rewardchain
rewardchaind keys list
rewardchaind keys show validator alice --keyring-backend file

 List all keys
rewardchaind keys list --keyring-backend os

# Get validator address
rewardchaind keys show validator -a --keyring-backend os

# Check balance (replace with your address)
rewardchaind query bank balances reward1abc123... --chain-id rewardchain

# Check specific token balance
rewardchaind query bank balance reward1abc123... ureward --chain-id rewardchain
