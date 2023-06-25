-- +goose Up
-- +goose StatementBegin

CREATE TABLE `article` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `author` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `title` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `body` text COLLATE utf8mb4_unicode_ci NOT NULL,
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` timestamp NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    FULLTEXT INDEX `author` (`author`),
    FULLTEXT INDEX `keyword` (`title`, `body`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `article`;
-- +goose StatementEnd
