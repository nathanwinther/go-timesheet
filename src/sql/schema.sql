CREATE TABLE IF NOT EXISTS `config` (
    `id` INTEGER PRIMARY KEY NOT NULL,
    `key` TEXT NOT NULL,
    `value` TEXT NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS `config_id1` ON `config`(`key`);

CREATE TABLE IF NOT EXISTS `user` (
    `id` INTEGER PRIMARY KEY NOT NULL,
    `key` TEXT NOT NULL,
    `active` INTEGER NOT NULL,
    `username` TEXT UNIQUE NOT NULL,
    `email` TEXT UNIQUE NOT NULL,
    `password` TEXT NOT NULL,
    `fullname` TEXT,
    `modified_date` INTEGER NOT NULL,
    `create_date` INTEGER NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS `user_ix1` ON `user`(`key`);
CREATE UNIQUE INDEX IF NOT EXISTS `user_ix2` ON `user`(`username`);
CREATE UNIQUE INDEX IF NOT EXISTS `user_ix3` ON `user`(`email`);

CREATE TABLE IF NOT EXISTS `user_client` (
    `id` INTEGER PRIMARY KEY NOT NULL,
    `user_id` INTEGER NOT NULL,
    `name` TEXT NOT NULL,
    `data` TEXT NOT NULL,
    `modified_date` INTEGER NOT NULL,
    `create_date` INTEGER NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS `user_client_ix1` ON `user_client`(`user_id`, `name`);

CREATE TABLE IF NOT EXISTS `user_session` (
    `id` INTEGER PRIMARY KEY NOT NULL,
    `key` TEXT NOT NULL,
    `user_id` INTEGER NOT NULL,
    `valid_until` INTEGER NOT NULL,
    `modified_date` INTEGER NOT NULL,
    `create_date` INTEGER NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS `user_session_ix1` ON `user_session`(`key`);
CREATE INDEX IF NOT EXISTS `user_session_ix2` ON `user_session`(`user_id`);

CREATE TABLE IF NOT EXISTS `user_verify` (
    `id` INTEGER PRIMARY KEY NOT NULL,
    `key` TEXT NOT NULL,
    `user_id` INTEGER NOT NULL,
    `valid_until` INTEGER NOT NULL,
    `modified_date` INTEGER NOT NULL,
    `create_date` INTEGER NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS `user_verify_ix1` ON `user_verify`(`key`);
CREATE INDEX IF NOT EXISTS `user_verify_ix2` ON `user_verify`(`user_id`);

