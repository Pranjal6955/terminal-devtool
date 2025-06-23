#!/bin/bash

# Test script for the backend API
# This script verifies all endpoints and functionality

# Set colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}===== Terminal DevTool Backend Test Script =====${NC}"

# Make sure we're in the backend directory
cd "$(dirname "$0")"

# Function to test an API endpoint
test_endpoint() {
  local name=$1
  local method=$2
  local endpoint=$3
  local data=$4
  local expected_status=$5
  
  echo -e "${YELLOW}Testing ${name}...${NC}"
  
  if [ "$method" == "GET" ]; then
    response=$(curl -s -w "\n%{http_code}" -X GET "http://localhost:8080$endpoint")
  else
    response=$(curl -s -w "\n%{http_code}" -X "$method" -H "Content-Type: application/json" -d "$data" "http://localhost:8080$endpoint")
  fi
  
  status=$(echo "$response" | tail -n1)
  body=$(echo "$response" | head -n -1)
  
  # Check if status matches expected
  if [ "$status" -eq "$expected_status" ]; then
    echo -e "${GREEN}✓ Status code: $status${NC}"
    echo "Response:"
    echo "$body" | jq . 2>/dev/null || echo "$body"
    return 0
  else
    echo -e "${RED}✗ Unexpected status code: $status (expected: $expected_status)${NC}"
    echo "Response:"
    echo "$body" | jq . 2>/dev/null || echo "$body"
    return 1
  fi
}

# 1. Check if the server is running
echo -e "${YELLOW}Checking if the server is running...${NC}"
if ! curl -s -o /dev/null http://localhost:8080/health; then
  echo -e "${RED}Server is not running. Please start the server first.${NC}"
  exit 1
fi

# 2. Test health check endpoint
test_endpoint "Health Check" "GET" "/health" "" 200

# 3. Create a test video if it doesn't exist
if [ ! -f "test_input.mp4" ]; then
  echo -e "${YELLOW}Creating test video...${NC}"
  ffmpeg -y -f lavfi -i "testsrc=duration=5:size=1280x720:rate=30" -c:v libx264 test_input.mp4
fi

# 4. Test media info endpoint
test_endpoint "Media Info" "GET" "/api/info?path=test_input.mp4" "" 200

# 5. Test process endpoint (MP4 to WebM)
cat > request.json << EOF
{
  "input": "test_input.mp4",
  "output": "test_output.webm",
  "resolution": "640x360",
  "bitrate": "500k",
  "format": "webm",
  "codec": "libvpx-vp9"
}
EOF

test_endpoint "Process Media (MP4 to WebM)" "POST" "/api/process" "$(cat request.json)" 200

# 6. Test compare endpoint
if [ -f "test_output.webm" ]; then
  cat > compare_request.json << EOF
  {
    "original": "test_input.mp4",
    "processed": "test_output.webm"
  }
EOF

  test_endpoint "Compare Media" "POST" "/api/compare" "$(cat compare_request.json)" 200
fi

# 7. Test process endpoint (MP4 to GIF)
cat > gif_request.json << EOF
{
  "input": "test_input.mp4",
  "output": "test_output.gif",
  "resolution": "320x240",
  "frame_rate": "10"
}
EOF

test_endpoint "Process Media (MP4 to GIF)" "POST" "/api/process" "$(cat gif_request.json)" 200

echo -e "${BLUE}===== Test Complete =====${NC}"
