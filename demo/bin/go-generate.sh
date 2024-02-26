#!/usr/bin/env bash
root_dir=$(cd "$(dirname "$0")"; cd ..; pwd)

protoExec=$(which "protoc")
if [ -z $protoExec ]; then
    echo 'Please install protoc!'
    exit 0
fi

protos_dir=$root_dir/protos
pb_dir=$root_dir/pb

mkdir -p $pb_dir

#delete old pb code.
rm -rf $root_dir/pb/*

echo "generating code"

echo "generating golang stubs..."
cd $protos_dir


# Generate OpenAPI JSON file
protoc -I $protos_dir --openapiv2_out $root_dir \
    --openapiv2_opt logtostderr=true \
    --openapiv2_opt=json_names_for_fields=false \
    $protos_dir/*.proto

# go grpc code
protoc -I $protos_dir \
    --go_out $root_dir/pb --go_opt paths=source_relative \
    --go-grpc_out $root_dir/pb --go-grpc_opt paths=source_relative \
    $protos_dir/*.proto

# http gw code
protoc -I $protos_dir --grpc-gateway_out $root_dir/pb \
    --grpc-gateway_opt logtostderr=true \
    --grpc-gateway_opt paths=source_relative \
    $protos_dir/*.proto

# cp golang client code
mkdir -p $root_dir/clients/go/pb

cp -R $root_dir/pb/*.go $root_dir/clients/go/pb

echo "generating golang code success"

echo "done!!!!"

exit 0
