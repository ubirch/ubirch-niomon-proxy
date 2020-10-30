# ubirch-niomon-proxy
Simple Ubirch Niomon Proxy

## usage

curl localhost:3000

curl -X POST localhost:3000/ubproxy/api/v1/upp

curl -s -H 'x-token: ea422f57-5e0b-4616-81bd-211ff49d21ce' -H 'x-ubirch-auth-type: ubirch'  -H 'x-ubirch-hardware-id: 82100f6f-ae55-44d0-a863-491a572d4569' -H 'x-ubirch-credential: MzVjYjEyMGYtMDRlZS00ZTViLWE1MGEtYjQ5ODFjNTE0YmE1' -H 'Content-Type: application/octet-stream' -X POST --data-binary @slack_upp.bin http://localhost:3000/ubproxy/api/v1/upp
