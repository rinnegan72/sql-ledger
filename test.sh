#!/bin/zsh

API_URL="http://localhost:8080"
USER_ID=72
CONCURRENT_WITHDRAWALS=15
CONCURRENT_DEPOSITS=15
WITHDRAWAL_AMOUNT=675
DEPOSIT_AMOUNT=500

RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

echo "=== Mixed Concurrent Transactions Test ==="
echo "Withdrawals: $CONCURRENT_WITHDRAWALS x $WITHDRAWAL_AMOUNT"
echo "Deposits: $CONCURRENT_DEPOSITS x $DEPOSIT_AMOUNT"
echo ""

TEMP_DIR=$(mktemp -d)
trap "rm -rf $TEMP_DIR" EXIT

# Launch withdrawals
for i in $(seq 1 $CONCURRENT_WITHDRAWALS); do
    {
        curl -s -X POST \
            -H "Content-Type: application/json" \
            -d "{\"amount\": $WITHDRAWAL_AMOUNT}" \
            "$API_URL/user/$USER_ID/withdraw" \
            -w "\n%{http_code}" \
            > "$TEMP_DIR/withdraw_$i.txt"
    } &
done

# Launch deposits
for i in $(seq 1 $CONCURRENT_DEPOSITS); do
    {
        curl -s -X POST \
            -H "Content-Type: application/json" \
            -d "{\"amount\": $DEPOSIT_AMOUNT}" \
            "$API_URL/user/$USER_ID/deposit" \
            -w "\n%{http_code}" \
            > "$TEMP_DIR/deposit_$i.txt"
    } &
done

wait

# Count results
WITHDRAW_SUCCESS=$(grep -l "200" "$TEMP_DIR"/withdraw_*.txt | wc -l)
DEPOSIT_SUCCESS=$(grep -l "200" "$TEMP_DIR"/deposit_*.txt | wc -l)

echo ""
echo -e "${GREEN}Successful withdrawals: $WITHDRAW_SUCCESS${NC}"
echo -e "${BLUE}Successful deposits: $DEPOSIT_SUCCESS${NC}"
echo ""
echo "Sample responses:"
echo -e "${RED}Withdrawal:${NC}"
head -1 "$TEMP_DIR/withdraw_1.txt"
echo -e "${BLUE}Deposit:${NC}"
head -1 "$TEMP_DIR/deposit_1.txt"