# RPC client and a server with Proof-of-Work anti-DDoS mechanics

## Proof of work principle

- Server generates 16-bytes sized Block
- Block includes a Salt and a Nonce (both 8-bytes sized)
- Block is being hashed using SHA-256 algo (32-bytes sized sequence as a result)
- Salt and a result of hashing are being sent to the Client as a puzzle
- Client must solve this puzzle and find Nonce calculating SHA-256 hashes from the known Salt and sequentially chosen number
- If puzzle solved, received and validated by the Server - Client will get a Wisdom Quote

![anti DDoS PoW](rpc_scheme.jpg)

## FAQ:

Q: What is Proof of Work how it linked to DDoS problem?
A: The answer is [here](https://en.wikipedia.org/wiki/Proof_of_work)

Q: Why do we need a Salt?
A: To prevent [Rainbow tables attack](https://en.wikipedia.org/wiki/Rainbow_table)

## Execution

Just say `make` to download all the dependencies, run tests and build a docker container for both client and server.

To run server:  
`$ make server`

To run client:  
`$ make client`
