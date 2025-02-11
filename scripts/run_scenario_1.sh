#!/bin/bash

echo "============================================================================"
echo " FUTURE TAKE HOME ASSIGNMENT"
echo "----------------------------------------------------------------------------"
echo "  * Demostrating the API service using curl"
echo 

# Set base URL
BASE_URL="http://localhost:8080/api/v1"

# Function to print test case
print_test() {
    echo 
    echo "~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~"
    echo "** TEST CASE: $1"
    echo "** EXPECTED: $2"
}

# Test Case 1: List appointments for trainer 1 (should be empty)
print_test "List appointments for trainer 1" "Empty list"
echo curl -s -w "\nStatus code: %{http_code}\n" "${BASE_URL}/appointments/trainers/1"
curl -s -w "\nStatus code: %{http_code}\n" "${BASE_URL}/appointments/trainers/1"

# Test Case 2: Get Availability between June 1 with small time range
print_test "Get Availability of trainer 1 on June 1" "List of all time slots on June 1 between 10AM and 2PM in UTC"
echo curl -s -w "\nStatus code: %{http_code}\n" \
    "${BASE_URL}/appointments/trainers/1/availability?starts_at=2025-06-01T18:00:00Z&ends_at=2025-06-01T22:00:00Z"
response="$(curl -s -w "\nStatus code: %{http_code}\n" \
    "${BASE_URL}/appointments/trainers/1/availability?starts_at=2025-06-01T18:00:00Z&ends_at=2025-06-01T22:00:00Z")"
echo "$response" | head -1 | jq .
echo "$response" | tail -1
echo -e "\n"

# Test Case 3: Try to book appointment outside business hours (6AM Pacific)
print_test "Book appointment outside business hours" "4xx ERROR - outside business hours"
echo curl -s -w "\nStatus code: %{http_code}\n" -X POST \
    "${BASE_URL}/appointments" \
    -H 'Content-Type: application/json' \
    -d '{
        "start_time": "2025-06-01T14:00:00Z",
        "end_time": "2025-06-01T14:30:00Z",
        "trainer_id": 1,
        "user_id": 10
    }'
response="$(curl -s -w "\nStatus code: %{http_code}\n" -X POST \
    "${BASE_URL}/appointments" \
    -H 'Content-Type: application/json' \
    -d '{
        "start_time": "2025-06-01T14:00:00Z",
        "end_time": "2025-06-01T14:30:00Z",
        "trainer_id": 1,
        "user_id": 10
    }')"
echo "$response" | head -1 | jq .
echo "$response" | tail -1
echo -e "\n"

# Test Case 4: Book valid appointment (11AM Pacific)
print_test "Book valid appointment at 11AM" "200 OK with new appointment"
echo curl -s -w "\nStatus code: %{http_code}\n" -X POST \
    "${BASE_URL}/appointments" \
    -H 'Content-Type: application/json' \
    -d '{
        "start_time": "2025-06-01T19:00:00Z",
        "end_time": "2025-06-01T19:30:00Z",
        "trainer_id": 1,
        "user_id": 10
    }'
response="$(curl -s -w "\nStatus code: %{http_code}\n" -X POST \
    "${BASE_URL}/appointments" \
    -H 'Content-Type: application/json' \
    -d '{
        "start_time": "2025-06-01T19:00:00Z",
        "end_time": "2025-06-01T19:30:00Z",
        "trainer_id": 1,
        "user_id": 10
    }')"
echo "$response" | head -1 | jq .
echo "$response" | tail -1
echo -e "\n"


# Test Case 5: Try to book conflicting appointment
print_test "Book conflicting appointment at 11AM" "4xx ERROR - time slot taken"
echo curl -s -w "\nStatus code: %{http_code}\n" -X POST \
    "${BASE_URL}/appointments" \
    -H 'Content-Type: application/json' \
    -d '{
        "start_time": "2025-06-01T19:00:00Z",
        "end_time": "2025-06-01T19:30:00Z",
        "trainer_id": 1,
        "user_id": 10
    }'
response="$(curl -s -w "\nStatus code: %{http_code}\n" -X POST \
    "${BASE_URL}/appointments" \
    -H 'Content-Type: application/json' \
    -d '{
        "start_time": "2025-06-01T19:00:00Z",
        "end_time": "2025-06-01T19:30:00Z",
        "trainer_id": 1,
        "user_id": 10
    }')"
echo "$response" | head -1 | jq .
echo "$response" | tail -1
echo -e "\n"

# Test Case 6: Book another valid appointment June 1:00PM
print_test "Book valid appointment for June 1 at 1PM" "200 OK with new appointment"
echo curl -s -w "\nStatus code: %{http_code}\n" -X POST \
    "${BASE_URL}/appointments" \
    -H 'Content-Type: application/json' \
    -d '{
        "start_time": "2025-06-01T21:00:00Z",
        "end_time": "2025-06-01T21:30:00Z",
        "trainer_id": 1,
        "user_id": 12
    }'
response="$(curl -s -w "\nStatus code: %{http_code}\n" -X POST \
    "${BASE_URL}/appointments" \
    -H 'Content-Type: application/json' \
    -d '{
        "start_time": "2025-06-01T21:00:00Z",
        "end_time": "2025-06-01T21:30:00Z",
        "trainer_id": 1,
        "user_id": 12
    }')"
echo "$response" | head -1 | jq .
echo "$response" | tail -1
echo -e "\n"

# Test Case 7: Check final availability
print_test "Get Availability of trainer 1 on June 1" "List of all time slots on June 1 between 10AM and 2PM in UTC, should not show 11AM and 1PM"
echo curl -s -w "\nStatus code: %{http_code}\n" \
    "${BASE_URL}/appointments/trainers/1/availability?starts_at=2025-06-01T18:00:00Z&ends_at=2025-06-01T22:00:00Z"
response="$(curl -s -w "\nStatus code: %{http_code}\n" \
    "${BASE_URL}/appointments/trainers/1/availability?starts_at=2025-06-01T18:00:00Z&ends_at=2025-06-01T22:00:00Z")"
echo "$response" | head -1 | jq .
echo "$response" | tail -1
echo -e "\n"

# Test Case 8: Get list of appointments for trainer 1 again
print_test "List appointments for trainer 1" "appointments for 11am and 1pm"
echo curl -s -w "\nStatus code: %{http_code}\n" "${BASE_URL}/appointments/trainers/1"
response="$(curl -s -w "\nStatus code: %{http_code}\n" "${BASE_URL}/appointments/trainers/1")"
echo "$response" | head -1 | jq .
echo "$response" | tail -1
echo
echo 
echo "End of API Tests"
