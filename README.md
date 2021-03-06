# Word Jumble

A little command-line tool to solve word games like [New York Times _Spelling Bee_](https://www.nytimes.com/puzzles/spelling-bee). By default, the reuse of letters is permitted. The `-c` (`--consume`) flag disables letter reuse, making it suitable for [Scrabble](https://scrabble.hasbro.com/en-us) or similar games.

## Usage

### Test one or more strings

```bash
wordjumble \
  [-v|--verbose] \
  [-c|--consume] \
  [-d|--dict <dictionary>] \
  [--use-array] \
  string_1 [string_2] [string_...] [string_n]
```

This will launch wordjumble loading the specified dictionary (or using the default `2of12inf`). Each string will be checked against the dictionary, separated by `------`. The `--dict` parameter can be used with any of the dictionaries from the `list` command (see below). The `--consume` parameter disables letter reuse unless the letter is explicitly duplicated in the input. The `--use-array` parameter tells the application to use the `arraytrie` implementation instead of the `maptrie` implementation.

#### Example

```bash
$ ./wordjumble abc def
baa
cab
------
deed
deeded
def
ed
fed
fee
feed
```

```bash
$ ./wordjumble -c abc def
cab
------
def
ed
```

### REPL mode

If no strings are provided, the application will present you with a prompt, `word>`. This is REPL (read-evaluate-print loop) mode. Strings can be entered here and they will be processed as normal. Exit the REPL with `^C` or `enter` on an empty line.

```bash
./wordjumble
Loaded dictionary 2of12inf: 81,883 words
words> abc
baa
cab
words> def
deed
deeded
def
ed
fed
fee
feed
words>
```

### List internal dictionaries

```bash
wordjumble list
```

#### Example (List)

```bash
$ ./wordjumble list
2of12
2of12inf
3esl
6of12
scrabble_words
words_alpha
```

## Implementation

There are two implementations in this application, both based on a [trie](https://en.wikipedia.org/wiki/Trie) structure -- An array-based trie or a map-based trie. For only unaccented Latin characters, the array-based true should be faster, but if you want to support dictionaries with large character sets, the array can take multiple steps per character. This is where the map-based implementation is possibly more effective. I haven't put any non-Latin dictionaries into the code, but the reason for this little toy was to play around with exactly that.

These trie implementations have another problem: false positives. Because each rune is split into bytes, the bytes of a particular rune are not held together. That means that the the implementation can return incorrect results. The most correct form would be a map of tries based on runes, though that would have a different interface.

## Dictionaries

The effectiveness of this kind of tool is tied to how well the dictionary matches the dictionary in the game. To that end, I've included a couple of different dictionaries.

* Several dictionaries from [aspell's 12 dicts](http://wordlist.aspell.net/12dicts-readme/)
  * 2of12
  * 2of12inf
  * 3esl
  * 6of12
* scrabble_words from [Collins Scrabble Words](https://drive.google.com/open?id=1oGDf1wjWp5RF_X9C7HoedhIWMh5uJs8s)
* words_alpha from [dwyl/english-words](dwyl/english-words)
