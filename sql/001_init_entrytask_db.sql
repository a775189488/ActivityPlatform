CREATE DATABASE IF NOT EXISTS entrytask_activity_platform_db 
	DEFAULT CHARACTER SET utf8mb4
	DEFAULT COLLATE utf8mb4_unicode_ci;

USE entrytask_activity_platform_db;

CREATE TABLE IF NOT EXISTS user_tab (
	id bigint(20) unsigned PRIMARY KEY AUTO_INCREMENT,
	aliasname varchar(64) NOT NULL,
	username varchar(64) NOT NULL UNIQUE,
	password varchar(255) NOT NULL,
	email varchar(255) NOT NULL,
	headpic varchar(255) NOT NULL,     # user head portrait address
	role tinyint NOT NULL,             # user role: 0 normal, 1 admin
	is_online bool NOT NULL DEFAULT 0,
	is_enable bool NOT NULL DEFAULT 1,
	create_time int unsigned NOT NULL,
	update_time int unsigned NOT NULL,
	delete_time int unsigned NOT NULL,
    INDEX idx_username (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS act_tab (
	id bigint(20) unsigned PRIMARY KEY AUTO_INCREMENT,
	title varchar(128) NOT NULL,
	begin_at int unsigned NOT NULL,
	end_at int unsigned NOT NULL,
	description varchar(255) NOT NULL,
	creator bigint(20) unsigned NOT NULL,
	act_type bigint(20) unsigned NOT NULL,
	address varchar(64) NOT NULL,
	status tinyint NOT NULL,
	create_time int unsigned NOT NULL,
	update_time int unsigned NOT NULL,
	INDEX idx_begin_at_end_at_type (`begin_at`,`end_at`,`act_type`),
	INDEX idx_end_at_type (`end_at`,`act_type`), 
	INDEX idx_type_begin_at (`act_type`,`begin_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS comment_tab (
	id bigint(20) unsigned PRIMARY KEY AUTO_INCREMENT,
	act_id bigint(20) unsigned NOT NULL,
	user_id bigint(20) unsigned NOT NULL,
	message varchar(255) NOT NULL,
	parent bigint(20) unsigned NOT NULL,
	create_time int unsigned NOT NULL,
	update_time int unsigned NOT NULL,
    INDEX idx_act_id (`act_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

# use for recording the relationship of user and activity
CREATE TABLE IF NOT EXISTS act_user_tab (
	id bigint(20) unsigned PRIMARY KEY AUTO_INCREMENT,
	act_id bigint(20) unsigned NOT NULL,   # activity id
	user_id bigint(20) unsigned NOT NULL,  # user id
	create_time int unsigned NOT NULL,
	update_time int unsigned NOT NULL,
	INDEX idx_user_id (`user_id`),
	INDEX idx_act_id (`act_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS act_type_tab (
	id bigint(20) unsigned PRIMARY KEY AUTO_INCREMENT,
	name varchar(64) NOT NULL,
	parent bigint(20) unsigned NOT NULL, # parent activity type id
	create_time int unsigned NOT NULL,
	update_time int unsigned NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;