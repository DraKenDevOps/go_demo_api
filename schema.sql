CREATE TABLE users
(
    user_id     INT AUTO_INCREMENT
        PRIMARY KEY,
    username    VARCHAR(50)                                                         NULL,
    telephone   VARCHAR(50)                                                         NULL,
    email       VARCHAR(50)                                                         NULL,
    password    VARCHAR(255)                                                        NULL,
    level       ENUM ('ADMIN', 'SUPPORT')                 DEFAULT 'SUPPORT'         NULL,
    role_action ENUM ('ALL', 'INSERT', 'UPDATE', 'QUERY') DEFAULT 'QUERY'           NULL,
    status      ENUM ('ACTIVE', 'DISABLE', 'SUSPEND')     DEFAULT 'ACTIVE'          NULL,
    op_id       INT                                                                 NULL,
    created_at  DATETIME                                  DEFAULT CURRENT_TIMESTAMP NOT NULL,
    -- updated_at  DATETIME                                                            NULL ON UPDATE CURRENT_TIMESTAMP,
    updated_at  DATETIME                                                            NULL,
    deleted_at  DATETIME                                                            NULL
);
CREATE INDEX users_index
    ON users (level, role_action, status, op_id, deleted_at);
