
#!/bin/bash
# for linux and macos

if [ -z "$1" ]; then
    echo "Usage: ./run.sh <PORT>"
    exit 1
fi

export PORT=$1
docker compose up -d --build
