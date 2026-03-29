cd ..
if [ -f ./bin/scut ]; then
    rm ./bin/scut
fi
go build -o ./bin/scut
./bin/scut