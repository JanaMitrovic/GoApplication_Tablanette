package constants

const (
	NEW_DECK_URL             = "https://www.deckofcardsapi.com/api/deck/new/"
	LIST_PILE_CARDS_URL      = "https://www.deckofcardsapi.com/api/deck/%s/pile/%s/list/"
	ADD_TO_PILE_URL          = "https://www.deckofcardsapi.com/api/deck/%s/pile/%s/add/?cards=%s"
	DRAW_CARDS_FROM_PILE_URL = "https://www.deckofcardsapi.com/api/deck/%s/pile/%s/draw/?cards=%s"
	NEW_SHUFFLED_DECK        = "https://www.deckofcardsapi.com/api/deck/new/shuffle/?deck_count=1"
	DRAW_A_CARD_URL          = "https://deckofcardsapi.com/api/deck/%s/draw/?count=%s"
	RETURN_TO_DECK_URL       = "https://www.deckofcardsapi.com/api/deck/%s/return/"
	AUTH_HEADER_MISSING      = "Auth header missing"
	INVALID_TOKEN            = "Invalid token"
	TOKEN_EXPIRED            = "Token expired"
	FORBIDDEN_ACCESS         = "Forbidden"
)
