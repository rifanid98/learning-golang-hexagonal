-- +goose Up
-- SQL in this section is executed when the migration is applied.
INSERT INTO `user` (`id`, `email`, `password`, `role`, `created_at`, `updated_at`) VALUES (555,'a@a.id', '$2a$08$cwJLLr.LfnKjUdpW6C3kE.KnEhdcQiVXGXcMd3iAXJ9IgMgLDnOci', 'admin', '2021-11-21 00:00:00', '2021-11-21 00:00:00');

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DELETE FROM `user` WHERE `id` = 555;
