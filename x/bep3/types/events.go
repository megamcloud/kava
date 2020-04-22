package types

// Events for bep3 module
const (
	EventTypeCreateAtomicSwap = "create_atomic_swap"
	EventTypeClaimAtomicSwap  = "claim_atomic_swap"
	EventTypeRefundAtomicSwap = "refund_atomic_swap"
	EventTypeSwapExpired      = "swap_expired"

	AttributeValueCategory       = ModuleName
	AttributeKeySender           = "sender"
	AttributeKeyRecipient        = "recipient"
	AttributeKeyAtomicSwapID     = "atomic_swap_id"
	AttributeKeyRandomNumberHash = "random_number_hash"
	AttributeKeyTimestamp        = "timestamp"
	AttributeKeySenderOtherChain = "sender_other_chain"
	AttributeKeyExpireHeight     = "expire_height"
	AttributeKeyAmount           = "amount"
	AttributeKeyExpectedIncome   = "expected_income"
	AttributeKeyDirection        = "direction"
	AttributeKeyClaimSender      = "claim_sender"
	AttributeKeyRandomNumber     = "random_number"
	AttributeKeyRefundSender     = "refund_sender"
	AttributeExpirationBlock     = "expiration_block"
)
