import os

directory = 'proto/'
output = "pkg/pb"
for proto_file in os.listdir(directory):
    if proto_file.endswith('.proto'):
        proto = directory + proto_file
        pb = proto_file.split('.')[0]
        command = f"protoc  \\\n    --go_out=. \\\n    --go_opt=M{proto}={output}/{pb}pb \\\n    --go-grpc_out=. \\\n    --go-grpc_opt=M{proto}={output}/{pb}pb \\\n    {proto}"
        print(command)