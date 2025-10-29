-- Initial schema creation (MySQL)

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id bigint unsigned NOT NULL AUTO_INCREMENT,
    username varchar(50) NOT NULL,
    password varchar(100) NOT NULL,
    email varchar(100) DEFAULT NULL,
    created_at timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY idx_username (username),
    UNIQUE KEY idx_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Tools table
CREATE TABLE IF NOT EXISTS tools (
    id bigint unsigned NOT NULL AUTO_INCREMENT,
    name varchar(100) NOT NULL,
    description text,
    icon varchar(255) DEFAULT NULL,
    plugin_name varchar(100) NOT NULL,
    is_enabled tinyint(1) DEFAULT 1,
    created_at timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY idx_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Tool histories table
CREATE TABLE IF NOT EXISTS tool_histories (
    id bigint unsigned NOT NULL AUTO_INCREMENT,
    user_id bigint unsigned DEFAULT NULL,
    tool_id bigint unsigned DEFAULT NULL,
    used_at timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    params text,
    result text,
    PRIMARY KEY (id),
    KEY idx_user_id (user_id),
    KEY idx_tool_id (tool_id),
    CONSTRAINT fk_tool_history_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE SET NULL,
    CONSTRAINT fk_tool_history_tool FOREIGN KEY (tool_id) REFERENCES tools (id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Notes table
CREATE TABLE IF NOT EXISTS notes (
    id varchar(100) NOT NULL,
    user_id bigint unsigned NOT NULL,
    title varchar(255) NOT NULL,
    content text NOT NULL,
    created_time timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    updated_time timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY idx_user_id (user_id),
    KEY idx_title (title),
    KEY idx_created_time (created_time),
    CONSTRAINT fk_note_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Login histories table
CREATE TABLE IF NOT EXISTS login_histories (
    id bigint unsigned NOT NULL AUTO_INCREMENT,
    username varchar(50) NOT NULL,
    ip_address varchar(50) DEFAULT NULL,
    success tinyint(1) NOT NULL DEFAULT 0,
    message varchar(255) DEFAULT NULL,
    user_agent text,
    login_time timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY idx_username (username),
    KEY idx_login_time (login_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Note histories table
CREATE TABLE IF NOT EXISTS note_history (
    id bigint unsigned NOT NULL AUTO_INCREMENT,
    note_id varchar(100) NOT NULL,
    title varchar(255) NOT NULL,
    content text NOT NULL,
    created_at timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY idx_note_id (note_id),
    KEY idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Teams table
CREATE TABLE IF NOT EXISTS team (
    id bigint unsigned NOT NULL AUTO_INCREMENT,
    name varchar(100) NOT NULL,
    description text,
    owner_id bigint unsigned NOT NULL,
    tenant_id bigint unsigned NOT NULL,
    created_at timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY idx_tenant_team_name (tenant_id,name),
    KEY idx_owner_id (owner_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Team members table
CREATE TABLE IF NOT EXISTS team_member (
    id bigint unsigned NOT NULL AUTO_INCREMENT,
    team_id bigint unsigned NOT NULL,
    user_id bigint unsigned NOT NULL,
    role varchar(50) NOT NULL DEFAULT 'member',
    tenant_id bigint unsigned NOT NULL,
    created_at timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY idx_team_id (team_id),
    KEY idx_user_id (user_id),
    KEY idx_tenant_id (tenant_id),
    UNIQUE KEY idx_team_user (team_id,user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Audit logs table
CREATE TABLE IF NOT EXISTS audit_logs (
    id bigint unsigned NOT NULL AUTO_INCREMENT,
    user_id bigint unsigned DEFAULT NULL,
    username varchar(50) DEFAULT NULL,
    action varchar(100) DEFAULT NULL,
    resource_type varchar(100) DEFAULT NULL,
    resource_id varchar(100) DEFAULT NULL,
    old_value text,
    new_value text,
    ip_address varchar(50) DEFAULT NULL,
    user_agent text,
    tenant_id bigint unsigned DEFAULT NULL,
    created_at timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY idx_user_id (user_id),
    KEY idx_tenant_id (tenant_id),
    KEY idx_action (action),
    KEY idx_resource_type (resource_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;