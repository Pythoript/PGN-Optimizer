package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"regexp"
	"strconv"
	"io/ioutil"
	"path/filepath"
	"github.com/fatih/color"
	"github.com/jessevdk/go-flags"
)

type Game struct {
	Headers map[string]string
	Moves   string
}

func ParsePGN(content string, removeComments bool) []Game {
	var games []Game
	currentGame := Game{Headers: make(map[string]string)}
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "[") {
			if len(currentGame.Moves) > 0 {
				if removeComments {
					currentGame.Moves = RemoveCommentsFromMoves(currentGame.Moves)
				}
				games = append(games, currentGame)
				currentGame = Game{Headers: make(map[string]string)}
			}
			header := strings.Trim(line, "[]")
			parts := strings.SplitN(header, " ", 2)
			if len(parts) == 2 {
				currentGame.Headers[parts[0]] = strings.Trim(parts[1], "\"")
			}
		} else if line != "" {
			currentGame.Moves += line + " "
		}
	}

	if len(currentGame.Moves) > 0 || len(currentGame.Headers) > 0 {
		games = append(games, currentGame)
	}

	return games
}

func RemoveEmptyHeaders(game *Game) {
	re := regexp.MustCompile(`^[\s?]*$`)
	for key, value := range game.Headers {
		if re.MatchString(value) {
			delete(game.Headers, key)
		}
	}
}

func RemoveFields(game *Game, fieldsToRemove []string) {
	for _, field := range fieldsToRemove {
		delete(game.Headers, field)
	}
}

func SavePGN(filename string, games []Game) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for i, game := range games {
		for key, value := range game.Headers {
			fmt.Fprintf(writer, "[%s \"%s\"]\n", key, value)
		}
		fmt.Fprintln(writer, strings.TrimSpace(game.Moves))
		if i < len(games)-1 {
			fmt.Fprintln(writer)
		}
	}

	return writer.Flush()
}

func RemoveGame(games *[]Game, index int) {
	if index >= 0 && index < len(*games) {
		*games = append((*games)[:index], (*games)[index+1:]...)
	}
}

func RemoveGamesWithoutPlayers(games *[]Game) {
	re := regexp.MustCompile(`^[\s?]*$`)
	filteredGames := []Game{}
	for _, game := range *games {
		whitePlayer, hasWhite := game.Headers["White"]
		blackPlayer, hasBlack := game.Headers["Black"]

		if (!hasWhite && !hasBlack) || (re.MatchString(whitePlayer) && re.MatchString(blackPlayer)) {
			continue
		}
		filteredGames = append(filteredGames, game)
	}
	*games = filteredGames
}

func FilterGamesByEloRange(games []Game, eloBelow, eloAbove int, filterEloless bool) []Game {
	filteredGames := []Game{}
	for _, game := range games {
		eloStrWhite := game.Headers["WhiteElo"]
		eloStrBlack := game.Headers["BlackElo"]

		eloWhite, errWhite := strconv.Atoi(eloStrWhite)
		eloBlack, errBlack := strconv.Atoi(eloStrBlack)

		if (filterEloless && (errWhite != nil || errBlack != nil)) ||
			(eloBelow != 0 && ((errWhite == nil && eloWhite < eloBelow) || (errBlack == nil && eloBlack < eloBelow))) ||
			(eloAbove != 0 && ((errWhite == nil && eloWhite > eloAbove) || (errBlack == nil && eloBlack > eloAbove))) {
			continue
		}

		filteredGames = append(filteredGames, game)
	}
	return filteredGames
}

func FilterGamesByYearRange(games []Game, yearBefore, yearAfter int, filterYearless bool) []Game {
	filteredGames := []Game{}
	for _, game := range games {
		dateStr, exists := game.Headers["Date"]
		if !exists && filterYearless {
			continue
		}
		yearExtracted, err := extractYear(dateStr)
		if err != nil && filterYearless{
			continue
		}
		if (yearBefore == 0 || yearExtracted >= yearBefore) && (yearAfter == 0 || yearExtracted <= yearAfter) {
			filteredGames = append(filteredGames, game)
		}
	}
	return filteredGames
}

func extractYear(dateStr string) (int, error) {
	yearPart := strings.Split(dateStr, ".")[0]
	return strconv.Atoi(yearPart)
}

func RemoveCommentsFromMoves(moves string) string {
	re := regexp.MustCompile(`\{[^}]*\}`)
	return re.ReplaceAllString(moves, "")
}

func findPGNFiles(dir string) ([]string, error) {
	var pgnFiles []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(info.Name()) == ".pgn" {
			pgnFiles = append(pgnFiles, path)
		}
		return nil
	})
	return pgnFiles, err
}

func concatenatePGNFiles(files []string) (string, error) {
	var contentBuilder strings.Builder
	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return "", err
		}
		contentBuilder.Write(data)
		contentBuilder.WriteString("\n\n")
	}
	return contentBuilder.String(), nil
}

var opts struct {
	Args struct {
		InputFile  string `name:"<input file>" description:"Input file path"`
	} `positional-args:"yes" required:"yes"`
	OutputFile string `short:"o" long:"output" default:"output.pgn" description:"Output file"`
	KeepEmpty  bool   `short:"e" long:"keep-empty" description:"Keep headers that are empty or contain only question marks or whitespace"`
	RemoveFields []string `short:"r" long:"remove-field" description:"Comma-separated list of fields to remove"`
	FilterBefore int  `long:"filter-before" description:"Filter out games before the specified year"`
	FilterAfter int  `long:"filter-after" description:"Filter out games after the specified year"`
	FilterYearlessGames bool  `long:"filter-yearless-games" description:"Filter out games that do not have a year"`
	FilterEloBelow int `long:"filter-elo-below" description:"Filter out games with ELO below this value"`
	FilterEloAbove int `long:"filter-elo-above" description:"Filter out games with ELO above this value"`
	FilterElolessGames bool `long:"filter-eloless-games" description:"Filter out games that do not have an ELO"`
	KeepPlayerless bool `short:"k" long:"keep-playerless-games" description:"Keep games were both players are unknown"`
	RemoveComments bool `long:"remove-comments" description:"Remove comments from the game moves"`
	ConcatPath  bool `short:"c" long:"concat" description:"recursively search for PGN files and concatenate them"`
}

func main() {

	parser := flags.NewParser(&opts, flags.Default)
	_, err := parser.Parse()

	if err != nil {
		fmt.Println("[!] Error parsing flags:", err)
		parser.WriteHelp(os.Stdout)
		os.Exit(1)
	}

	green := color.New(color.FgGreen, color.Bold)
	red := color.New(color.FgRed, color.Bold)

	var content string
	if opts.ConcatPath {
		pgnFiles, err := findPGNFiles(opts.Args.InputFile)
		if err != nil {
			red.Printf("\n[!] Error finding PGN files: %v", err)
			os.Exit(1)
		}

		content, err = concatenatePGNFiles(pgnFiles)
		if err != nil {
			red.Printf("\n[!] Error concatenating PGN files: %v", err)
			os.Exit(1)
		}

		green.Printf("\n[+] Concatenated %d PGN files", len(pgnFiles))
	} else {
		data, err := ioutil.ReadFile(opts.Args.InputFile)
		if err != nil {
			red.Printf("\n[!] Error reading file: %v", err)
			os.Exit(1)
		}
		content = string(data)
	}

	games := ParsePGN(string(content), opts.RemoveComments)
	gameCount := len(games)
	green.Printf("\n[+] Parsed %d games\n", gameCount)

	fieldsToRemove := append([]string{"ECO", "PlyCount", "Variation"}, opts.RemoveFields...)
	if len(opts.RemoveFields) > 0 {
		for i := range games {
			RemoveFields(&games[i], fieldsToRemove)
		}
	}
	if !opts.KeepEmpty {
		for i := range games {
			RemoveEmptyHeaders(&games[i])
		}
	}

	if !opts.KeepPlayerless {
		RemoveGamesWithoutPlayers(&games)
		green.Printf("[+] Removed %d game(s) without Player names\n", gameCount-len(games))
	}

	if opts.FilterBefore != 0 || opts.FilterAfter != 0 || opts.FilterYearlessGames {
		games = FilterGamesByYearRange(games, opts.FilterBefore, opts.FilterAfter, opts.FilterYearlessGames)
		green.Printf("[+] Filtered out %d games based on year range.\n", gameCount-len(games))
	}

	gameCount = len(games)
	if opts.FilterEloBelow != 0 || opts.FilterEloAbove != 0 || opts.FilterElolessGames {
		games = FilterGamesByEloRange(games, opts.FilterEloBelow, opts.FilterEloAbove, opts.FilterElolessGames)
		green.Printf("[+] Filtered out %d games based on ELO range.\n", gameCount-len(games))
	}

	err = SavePGN(opts.OutputFile, games)
	if err != nil {
		red.Printf("[!] Error saving file: %v\n", err)
	}

	green.Printf("[+] Saved %d games to %s\n", len(games), opts.OutputFile)
}