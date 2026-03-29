from pathlib import Path


def gather_go_files(root_dir=".", output_file="all_go_files.txt"):
    root = Path(root_dir)
    go_files = sorted(root.rglob("*.go"))

    with open(output_file, "w", encoding="utf-8") as out:
        for file_path in go_files:
            rel_path = file_path.relative_to(root)
            out.write(f"{'=' * 80}\n")
            out.write(f"FILE: {rel_path}\n")
            out.write(f"{'=' * 80}\n\n")

            try:
                content = file_path.read_text(encoding="utf-8")
            except UnicodeDecodeError:
                content = file_path.read_text(encoding="utf-8", errors="replace")

            out.write(content)
            out.write("\n\n")

    print(f"Collected {len(go_files)} Go files into {output_file}")


if __name__ == "__main__":
    gather_go_files()
