# Transaction Ledger for MySQL
----

You can build a simple transaction ledger. 
Using MySQL. build it out as APIs
use DB transactions for locking and updating things. 
You can even try to add complexity 
by adding extensive logging and telemetry support

## Requirements

1. Transaction ledger
2. observabillity (optional)

## Logic

1. Create an API call for both withdrawing and despositing money
2. Implement Row locking to avoid Race conditions
where simultaneuos transaction affecting
the same row don't cause issues

## Pseudocode

```
fn main()
{
  while true{
    post {hostname}/user/{user-id}/deposit {
      result = update(user_id, amount, "desposit")
      response(result)
    }
    post {hostname}/user/{user-id}/withdraw {
      result = update(user_id, amount, "withdraw")
      response(result)
    }
  }
}

fn update(user_id, amount, type){
tx := db.begin_transaction() 
  try {
    current_balance := tx.get_balance_for_update(user_id)
    if (type == "withdraw" and current_balance < amount){
      tx.rollback()
      return "not enough funds"
    }
    transaction_type = type == withdraw ? -1 : 1
    update_amount = amount * transaction_type
    tx.update_balance(user_id, current_balance + update_amount )
    tx.commit()
    return "{type} completed successfully"
  }
  catch{
    tx.rollback()
    return "{type} failed try again"
  }
}
```

## Actual architecture

### Input Json with POST
withdrawl example:
```
curl -X POST \
-d '{ amount: 675 }' \
http://localhost:8080/user/72/withdraw
```

expected response:
```
{
  "status": "success",
  "message": "transaction completed successfully"
}
```
