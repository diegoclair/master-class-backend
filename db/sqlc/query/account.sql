-- name: CreateAccount :one
INSERT INTO accounts (
  owner, 
  balance, 
  currency
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts 
WHERE id = $1 LIMIT 1
-- FOR UPDATE;  -- this is necessary when we are working with multithread (if some thread started a transaction, it will wait the transaction to finish before do the select and it will get the actual value)
FOR NO KEY UPDATE; -- FOR UPDATE only do not solve our problema because we got deadlock when the account id is used to update the balance, it generates lock, so we tell to database that we won't update the account id (we just use as WHERE condition)

-- name: ListAccounts :many
SELECT * FROM accounts
ORDER BY id
-- pagination
LIMIT $1
OFFSET $2;

-- name: UpdateAccount :one
UPDATE accounts
SET balance = $2
WHERE id = $1
RETURNING *;

-- name: AddAccountbalance :one
UPDATE accounts
SET balance = balance + sqlc.arg(amount) -- we use sqlc.arg when we want that the parameter had a different name that from db field name (instead of balance it will be amount)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts WHERE id = $1;