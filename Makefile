# 1. ПСЕВДОЦЕЛИ (.PHONY)
# Говорят утилите make: "proto, db-up и db-down — это названия команд, а не реальные файлы на диске".
# Без этого, если бы у тебя в папке случайно появился файл с именем "proto", make подумал бы, что он уже создан, и не выполнял бы команду.
.PHONY: proto db-up db-down

# 2. ЦЕЛЬ: proto
# Запускается командой: make proto
proto:
	protoc --proto_path=pkg/proto \
	       --go_out=pkg/gen --go_opt=paths=source_relative \
	       --go-grpc_out=pkg/gen --go-grpc_opt=paths=source_relative \
	       pkg/proto/*.proto