-- Rollback initial schema

-- 删除所有表（按照依赖关系逆序）
DROP TABLE IF EXISTS note_history;
DROP TABLE IF EXISTS tool_histories;
DROP TABLE IF EXISTS login_histories;
DROP TABLE IF EXISTS notes;
DROP TABLE IF EXISTS tools;
DROP TABLE IF EXISTS users;