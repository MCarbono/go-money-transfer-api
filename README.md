




consultar o balance: 

curl localhost:3000/balance/1
1 - id do usuário


 curl --location 'http://localhost:3000/transfer' \
    --header 'Content-Type: application/json' \
    --data '{
        "amount": 100.0,
        "debtor_id": 1,
        "beneficiary_id": 2
    }'