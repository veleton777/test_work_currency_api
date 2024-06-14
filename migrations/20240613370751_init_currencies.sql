-- +goose Up
-- +goose StatementBegin
INSERT INTO currencies (id, name, code, type, is_available)
VALUES ('064b095f-3cb2-40cb-a4f9-6843148cfcc7', 'Dollar', 'USD', 2, true),
       ('3aa3c627-ae64-44dd-8048-61827d2ee3bc', 'Bitcoin', 'BTC', 1, true),
       ('24985ff1-4ecb-49a2-95ea-90ccaa28ace9', 'Ethereum', 'ETH', 1, true)
;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE
FROM currencies
WHERE id in (
             '064b095f-3cb2-40cb-a4f9-6843148cfcc7',
             '3aa3c627-ae64-44dd-8048-61827d2ee3bc',
             '24985ff1-4ecb-49a2-95ea-90ccaa28ace9'
    );
-- +goose StatementEnd
