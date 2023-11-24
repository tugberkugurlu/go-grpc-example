# go-grpc-example

Call the server:

```
@tugberkugurlu ➜ /workspaces/go-grpc-example (main) $ grpcurl -plaintext -d '{"name": "Tugberk"}'  localhost:50051 hellowo
rld.Greeter/SayHello
{rld.Greeter/SayHell
  "message": "Hello Tugberk"
}
```