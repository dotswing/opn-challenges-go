### This project is assignment for OPN from Nattapat Duangkaew

#### Command options
-requestPerSec: Number of request/sec, default to number cpu in the running machine

-file: Path to encrypted file ex: fng.1000.csv.rot128, '/Users/{user}/Desktop/fng.1000.csv.rot128'

#### Run code
```sh
go run cmd/go-tamboon/main.go -requestPerSec 4 -file fng.1000.csv.rot128
#OR
./bin/go-tamboon -requestPerSec 4 -file fng.1000.csv.rot128
```



#### Build Command
```sh
OMISE_PUBLIC_KEY="<omise_public_key>" OMISE_SECRET_KEY="<omise_secret_key>" go build -o ./bin/go-tamboon cmd/go-tamboon/main.go
```

#### Screenshot from finished result
There are a lot of errors due to low rate limit on Vault api for token creation

To make it pass rate limit we could specify `-requestPerSec` option to `1`

![My Image](images/finished.png "Optional Title")
