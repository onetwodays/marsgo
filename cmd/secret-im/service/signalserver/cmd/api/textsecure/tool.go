package textsecure
func GetEnvelopeType(typ int) Envelope_Type  {
	switch typ {
	case 0:
		return Envelope_UNKNOWN
	case 1:
		return Envelope_CIPHERTEXT
	case 2:
		return Envelope_KEY_EXCHANGE
	case 3:
		return Envelope_PREKEY_BUNDLE
	case 5:
		return Envelope_RECEIPT
	case 6:
		return Envelope_UNIDENTIFIED_SENDER
	default:
		return Envelope_UNKNOWN
	}
}
