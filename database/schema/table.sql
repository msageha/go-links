USE go_links;

CREATE TABLE session (
    token varchar(2048) NOT NULL,
    user_id bigint NOT NULL,
    created_at datetime NOT NULL default current_timestamp,
    updated_at datetime NOT NULL default current_timestamp on update current_timestamp,
    PRIMARY KEY(token(128), user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE user (
    id bigint UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    email varchar(128) UNIQUE NOT NULL,
    is_deleted boolean NOT NULL default 0,
    created_at datetime NOT NULL default current_timestamp,
    updated_at datetime NOT NULL default current_timestamp on update current_timestamp
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE short_link (
    id bigint UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id bigint UNSIGNED NOT NULL,
    shortpath varchar(200) NOT NULL,
    shortpath_prefix varchar(200) NOT NULL,
    destination_url varchar(3000) NOT NULL,
    visits_count bigint UNSIGNED NOT NULL default 0,
    visits_count_last_updated datetime,
    namespace varchar(30),
    type varchar(30),
    display_shortpath varchar(200),
    created_at datetime NOT NULL default current_timestamp,
    updated_at datetime NOT NULL default current_timestamp on update current_timestamp
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
