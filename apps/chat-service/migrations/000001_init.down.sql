DROP INDEX IF EXISTS idx_messages_chat_id_id;
DROP INDEX IF EXISTS idx_chat_members_user_id;
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS chats;
DROP TYPE IF EXISTS chat_type;
DROP TABLE IF EXISTS chat_members;