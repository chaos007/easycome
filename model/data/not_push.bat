set outdir=pb
::set protodir = ..\proto
::mkdir %outdir%
protoc.exe --go_out=plugins=grpc:%outdir% --proto_path "proto" proto\*.proto
::protoc.exe --csharp_out=client --proto_path "proto" proto\*.proto
protoc.exe --plugin=protoc-gen-meta=protoc-gen-msgcode.exe --proto_path "proto" --meta_out=:. proto\*.proto