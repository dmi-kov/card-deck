CREATE TABLE IF NOT EXISTS decks
(
    deck_id     UUID                     NOT NULL PRIMARY KEY,
    is_shuffled BOOL                     NOT NULL DEFAULT FALSE,
    remaining   INTEGER                  NOT NULL,
    card_codes  TEXT[]                   NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);