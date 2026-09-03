#!/usr/bin/env python3
"""
repo_to_txt.py — конвертер Go-репозитория в текст для LLM.

Использование:
    python3 repo_to_txt.py                  # текущая папка → repo.txt
    python3 repo_to_txt.py /path/to/project # указанная папка → repo.txt
    python3 repo_to_txt.py . output.txt     # указать имя выходного файла
"""

import os
import sys

# ── Настройки ────────────────────────────────────────────────────────────────

# Расширения файлов, которые включаем
INCLUDE_EXTENSIONS = {".go", ".mod", ".sum", ".sql", ".yaml", ".yml", ".env.example", ".Dockerfile", "Dockerfile", ".html", ".js", ".css", ".conf"}

# Папки, которые полностью пропускаем
SKIP_DIRS = {".git", ".idea", ".vscode", "vendor", "node_modules", "__pycache__"}

# Файлы, которые пропускаем
SKIP_FILES = {"repo.txt", "repo_to_txt.py"}

# Максимальный размер одного файла (байт). Большие файлы будут обрезаны.
MAX_FILE_SIZE = 100_000

# ── Логика ───────────────────────────────────────────────────────────────────

def should_include(path: str) -> bool:
    """Проверяет, нужно ли включать файл."""
    filename = os.path.basename(path)
    if filename in SKIP_FILES:
        return False
    _, ext = os.path.splitext(filename)
    # Dockerfile без расширения
    if filename == "Dockerfile":
        return True
    return ext in INCLUDE_EXTENSIONS


def build_tree(root: str) -> list[str]:
    """Строит список строк дерева директорий."""
    lines = []

    def walk(dir_path: str, prefix: str):
        try:
            entries = sorted(os.scandir(dir_path), key=lambda e: (e.is_file(), e.name))
        except PermissionError:
            return

        entries = [e for e in entries if not (e.is_dir() and e.name in SKIP_DIRS)]
        entries = [e for e in entries if not (e.is_file() and not should_include(e.path))]

        for i, entry in enumerate(entries):
            is_last = (i == len(entries) - 1)
            connector = "└── " if is_last else "├── "
            lines.append(f"{prefix}{connector}{entry.name}")
            if entry.is_dir():
                extension = "    " if is_last else "│   "
                walk(entry.path, prefix + extension)

    lines.append(f"└── {os.path.basename(root)}/")
    walk(root, "    ")
    return lines


def collect_files(root: str) -> list[str]:
    """Возвращает список путей к файлам в порядке обхода."""
    result = []
    for dirpath, dirnames, filenames in os.walk(root):
        # Пропускаем скрытые и системные папки (изменение на месте)
        dirnames[:] = sorted(d for d in dirnames if d not in SKIP_DIRS)
        for filename in sorted(filenames):
            filepath = os.path.join(dirpath, filename)
            if should_include(filepath):
                result.append(filepath)
    return result


def read_file(path: str) -> str:
    """Читает файл, обрезает если слишком большой."""
    size = os.path.getsize(path)
    try:
        with open(path, "r", encoding="utf-8") as f:
            if size > MAX_FILE_SIZE:
                content = f.read(MAX_FILE_SIZE)
                return content + f"\n\n... [файл обрезан, показано {MAX_FILE_SIZE} байт из {size}]"
            return f.read()
    except UnicodeDecodeError:
        return "[бинарный файл, содержимое пропущено]"


def main():
    # Аргументы командной строки
    root = sys.argv[1] if len(sys.argv) > 1 else "."
    output = sys.argv[2] if len(sys.argv) > 2 else "repo.txt"

    root = os.path.abspath(root)
    if not os.path.isdir(root):
        print(f"Ошибка: '{root}' не является директорией.")
        sys.exit(1)

    print(f"📂 Сканирую: {root}")

    sections = []

    # ── Дерево директорий ────────────────────────────────────────────────────
    tree_lines = build_tree(root)
    sections.append("## Directory Structure\n\n" + "\n".join(tree_lines))

    # ── Содержимое файлов ────────────────────────────────────────────────────
    files = collect_files(root)
    print(f"📄 Найдено файлов: {len(files)}")

    file_sections = []
    for filepath in files:
        rel = os.path.relpath(filepath, root)
        content = read_file(filepath)
        file_sections.append(f"## File: {rel}\n\n```\n{content}\n```")

    sections.append("\n\n---\n\n".join(file_sections))

    # ── Запись ───────────────────────────────────────────────────────────────
    final = "\n\n---\n\n".join(sections)
    with open(output, "w", encoding="utf-8") as f:
        f.write(final)

    size_kb = os.path.getsize(output) / 1024
    print(f"✅ Готово: {output} ({size_kb:.1f} KB)")


if __name__ == "__main__":
    main()
