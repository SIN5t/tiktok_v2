
#!/bin/bash
CURDIR=$(cd $(dirname $0); pwd)
BinaryName=api-gateway
echo "$CURDIR/bin/${BinaryName}"
exec $CURDIR/bin/${BinaryName}
