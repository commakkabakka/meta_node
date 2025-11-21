一、环境

-   Ubuntu 24.04.3 LTS
-   go1.25.4

二、运行

```
go mod tidy
go run .
```

三、编译智能合约，生成 ABI 和字节码文件.

```
solc --abi --bin --overwrite Counter.sol -o ./build
```

四、使用 abigen 生成 Go 绑定代码

```
mkdir counter
./abigen --abi=./build/Counter.abi --bin=./build/Counter.bin --pkg=counter --out=./counter/Counter.go
go mod tidy
```
