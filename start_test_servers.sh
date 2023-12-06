#! /usr/bin/env bash

# You can use this script to test conduit. It starts all test servers and waits for them to finish.
# Press Ctrl+C to stop all servers.

trap handle_interruption INT

pids=()

start_server () {
    python3 server.py "$1" "$2" 2> /dev/null > /dev/null &
    pids+=($!)
    sleep 0.5
}

handle_interruption () {
    echo "Stopping servers..."
    for pid in "${pids[@]}"; do
        kill $pid
    done
}

# Inicializando os servidores
echo "Initializing first server on port 3000..."
start_server "server_1" 3000

echo "Initializing second server on port 3001..."
start_server "server_2" 3001

echo "Initializing third server on port 3002..."
start_server "server_3" 3002

echo "Initializing fourth server on port 3003..."
start_server "server_4" 3003

echo "
** Running test servers... **
** Press Ctrl+C to stop **
"

# Aguarda todos os servidores terminarem
wait
