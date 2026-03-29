#!/bin/bash

# Configuration
BASE_URL="http://localhost:9820/v1/orders"

# Number of requests to make (default: 10)
NUM_REQUESTS=${1:-10}

# Delay between requests in seconds (default: 0.5)
DELAY=${2:-0.5}

echo "Starting order service simulation with $NUM_REQUESTS requests..."
echo "Delay between requests: ${DELAY}s"
echo "----------------------------------------"

# Function to create an order
create_order() {
    local pos_id=$1
    local price=$2
    local recipe_id=$3
    
    curl --location --request POST "${BASE_URL}/create" \
        --form "pos_id=${pos_id}" \
        --form "price=${price}" \
        --form "recipe_id=${recipe_id}" \
        --silent --show-error
}

# Function to get an order
get_order() {
    local order_id=$1
    
    curl --location --request GET "${BASE_URL}/${order_id}" \
        --silent --show-error
}

# Function to list orders
list_orders() {
    curl --location --request GET "${BASE_URL}/list" \
        --silent --show-error
}

# Function to update an order
update_order() {
    local order_id=$1
    local pos_id=$2
    local price=$3
    local recipe_id=$4
    
    curl --location --request PUT "${BASE_URL}/${order_id}" \
        --form "pos_id=${pos_id}" \
        --form "price=${price}" \
        --form "recipe_id=${recipe_id}" \
        --silent --show-error
}

# Simulate various operations
for i in $(seq 1 $NUM_REQUESTS); do
    echo "Request $i/$NUM_REQUESTS"
    
    # Generate variation in the data
    POS_ID=$((10000 + $i))
    PRICE=$((30000 + ($i * 5000)))
    RECIPE_ID=$((($i % 10) + 1))  # Cycle through recipe IDs 1-10
    
    # Perform different operations based on request number
    case $((i % 4)) in
        0)
            echo "Creating order: pos_id=${POS_ID}, price=${PRICE}, recipe_id=${RECIPE_ID}"
            create_order $POS_ID $PRICE $RECIPE_ID
            ;;
        1)
            echo "Getting order with ID: $i"
            get_order $i
            ;;
        2)
            echo "Listing orders"
            list_orders
            ;;
        3)
            echo "Updating order ID: $i with new price: $PRICE"
            update_order $i $POS_ID $PRICE $RECIPE_ID
            ;;
    esac
    
    echo ""
    echo "Response for request $i completed"
    
    # Add delay between requests (except for the last one)
    if [ $i -lt $NUM_REQUESTS ]; then
        echo "Waiting ${DELAY}s before next request..."
        sleep $DELAY
    fi
    
    echo "----------------------------------------"
done

echo ""
echo "Simulation completed!"
echo "Total requests made: $NUM_REQUESTS"