# PGN Parser and Optimizer

This is a command-line tool written in Go for parsing, filtering, and optimizing chess games in PGN format. It allows you to filter games by player ratings (ELO), date range, and other criteria, remove specific headers or comments, concatenate multiple PGN files, and output the processed games to a new PGN file.

## Features

- **Parse PGN Files:** Reads PGN files and extracts game headers and moves.
- **Remove Empty Headers:** Option to remove headers that are empty or contain only whitespace or question marks.
- **Remove Specific Fields:** Allows the removal of selected headers such as ECO, PlyCount, Variation, or any other custom fields.
- **Filter by Year and ELO:** Filter games based on ELO rating or date range, and exclude games without these values.
- **Remove Comments:** Strip out comments from game moves to reduce file size and simplify game records.
- **Concatenate PGN Files:** Option to recursively search a directory for PGN files, concatenate them, and process as a single file.

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/Pythoript/PGN-Optimizer.git
   ```
2. Install dependencies:
   ```bash
   go mod init PGN-Optimizer
   go get github.com/fatih/color
   go get github.com/jessevdk/go-flags
   ```
3. Build the project:
   ```bash
   go build index.go -o pgn-parser
   ```

## Usage

```bash
./pgn-parser [options] <input-file>
```

### Options

| Option                    | Description                                                                                          |
|---------------------------|------------------------------------------------------------------------------------------------------|
| `<input file>`            | Path to the input PGN file or directory if concatenating files.                                      |
| `-o, --output`            | Path for the output file. Default: `output.pgn`.                                                     |
| `-e, --keep-empty`        | Keep headers that are empty or contain only whitespace or question marks.                            |
| `-r, --remove-field`      | Comma-separated list of fields to remove from headers.                                               |
| `--filter-before`         | Filter out games played before the specified year.                                                   |
| `--filter-after`          | Filter out games played after the specified year.                                                    |
| `--filter-yearless-games` | Exclude games that do not contain a date header.                                                     |
| `--filter-elo-below`      | Exclude games with ELO below the specified value.                                                    |
| `--filter-elo-above`      | Exclude games with ELO above the specified value.                                                    |
| `--filter-eloless-games`  | Exclude games without ELO ratings.                                                                   |
| `-k, --keep-playerless-games` | Keep games where both players are unknown.                                                     |
| `--remove-comments`       | Remove comments from game moves.                                                                     |
| `-c, --concat`            | Recursively search a directory for PGN files and concatenate them.                                   |

### Example Commands

- **Parse and save PGN files with ELO filtering:**
  ```bash
  ./pgn-parser --filter-elo-below 1200 --filter-elo-above 2000 example.pgn
  ```
- **Remove empty headers and comments from moves:**
  ```bash
  ./pgn-parser --remove-comments --keep-empty=false example.pgn
  ```
- **Concatenate and parse all PGN files in a directory:**
  ```bash
  ./pgn-parser --concat --output all_games.pgn /path/to/pgn/files
  ```

## License

This project is licensed under the MIT License.

---

This tool is ideal for chess enthusiasts and developers working with large PGN files who want to filter, clean, and optimize game data. Contributions are welcome!
