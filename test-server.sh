#!/bin/bash

# Manual test script for Redis server
# Tests all features: SET, GET, DEL, KEYS, SIZE, FLUSH, PING

echo "🧪 Testing Redis Server on localhost:6378"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Function to send command and get response
test_command() {
    local cmd="$1"
    local expected="$2"
    echo "📤 Sending: $cmd"
    response=$(echo "$cmd" | nc localhost 6378 | head -1)
    echo "📥 Response: $response"
    
    if [ -n "$expected" ] && [ "$response" != "$expected" ]; then
        echo "❌ FAILED: Expected '$expected'"
    else
        echo "✅ OK"
    fi
    echo ""
}

echo "1️⃣  Testing PING"
test_command "PING" "+PONG"

echo "2️⃣  Testing SET"
test_command "SET user1 Alice" "+OK"

echo "3️⃣  Testing GET"
test_command "GET user1" "Alice"

echo "4️⃣  Testing SET with spaces in value"
test_command "SET greeting Hello World from Redis"
test_command "GET greeting" "Hello World from Redis"

echo "5️⃣  Testing Multiple SET operations"
test_command "SET user2 Bob"
test_command "SET user3 Charlie"
test_command "SET user4 David"

echo "6️⃣  Testing KEYS (list all keys)"
echo "📤 Sending: KEYS"
echo "📥 Response:"
echo "KEYS" | nc localhost 6378
echo "✅ Should show: user1, user2, user3, user4, greeting"
echo ""

echo "7️⃣  Testing SIZE"
test_command "SIZE" "5"

echo "8️⃣  Testing DEL"
test_command "DEL user2" "+OK"
test_command "GET user2" "(nil)"
test_command "SIZE" "4"

echo "9️⃣  Testing GET non-existent key"
test_command "GET nonexistent" "(nil)"

echo "🔟 Testing case insensitivity"
test_command "set lowercase value" "+OK"
test_command "GET lowercase" "value"
test_command "del lowercase" "+OK"

echo "1️⃣1️⃣  Testing error handling"
test_command "INVALID" "ERR unknown command 'INVALID'"
test_command "SET" "ERR wrong number of arguments for 'set' command"
test_command "GET" "ERR wrong number of arguments for 'get' command"

echo "1️⃣2️⃣  Testing FLUSH"
test_command "FLUSH" "+OK"
test_command "SIZE" "0"
test_command "KEYS" "(empty)"

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🎉 Manual tests complete!"
echo ""
echo "💡 You can also test manually:"
echo "   nc localhost 6378"
echo "   Then type commands interactively"

