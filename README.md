# Word Jumble

A little command-line tool to solve word games like [New York Times _Spelling Bee_](https://www.nytimes.com/puzzles/spelling-bee). This permits the reuse of letters, making it unsuitable for [Scrabble](https://scrabble.hasbro.com/en-us).

## Usage

### Test one or more strings

```bash
wordjumble [-v|--verbose] [-d|--dict <dictionary>] [--use-array] string_1 [string_2] [string_...] [string_n]
```

This will launch wordjumble loading the specified dictionary (or using the default `2of12inf`). Each string will be checked against the dictionary, separated by `------`. The `--dict` parameter can be used with any of the dictionaries from the `list` command (see below). The `--use-array` parameter tells the application to use the `arraytrie` implementation instead of the `maptrie` implementation.

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

## Dictionaries

The effectiveness of this kind of tool is tied to how well the dictionary matches the dictionary in the game. To that end, I've included a couple of different dictionaries.

* Several dictionaries from [aspell's 12 dicts](http://wordlist.aspell.net/12dicts-readme/)
  * 2of12
  * 2of12inf
  * 3esl
  * 6of12
* scrabble_words from [Collins Scrabble Words](https://drive.google.com/open?id=1oGDf1wjWp5RF_X9C7HoedhIWMh5uJs8s)
* words_alpha from [dwyl/english-words](dwyl/english-words)
