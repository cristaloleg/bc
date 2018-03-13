## BC - blockchain

To start:
```sh
go build main.go && ./main
```

and visit http://localhost:3456/

New block - http://localhost:3456/mine?data=myvalue
All blocks - http://localhost:3456/blocks
WS peers - http://localhost:3456/peers

To test:
```sh
go test ./...
```

coverage: 93.0% of statements

Use `config.json` to change config