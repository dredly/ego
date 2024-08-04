# ego
A simple command line tool for elo tracking written in go.
Uses a local sqlite database under the hood.

## Commands

### create
```ego create```

Initialises a local sqlite db in the location  ```~/.ego/ego.db```


### add
```ego add -name <player name>```

Adds a player to the leaderboard with the specified name and a default starting elo of 1000. Name must be unique.

```ego add -name <player name> -elo <elo>```

Adds a player to the leaderboard with the specified name and starting elo

### leaderboard
```ego leaderboard```

Displays the leaderboard

### games
```ego games```

Displays all previous games (newest to oldest)

```ego games -limit <n>```

Displays the ```n``` most recent games

```ego games -player <player name>```

Displays all previous games involving one player

### record
```ego record -p1 <name of player 1> -p2 <name of player 2> 11-3```

Records an 11-3 win for player 1 over player 2. Note that if the loser scores 0 points, the elo change will be doubled

### undo
```ego undo```

Undo recording the last game. Removes the game itself from the database, and sets the players ELOs back to what they were before the game

### stats

Available stats are: games played, win rate and peak ELO.
If no stats specified, all stats will be shown

```ego stats -name=Bob```

Shows all stats for Bob

```ego stats -name=Bob -peak```

Shows just the peak ELO for Bob

```ego stats -name=Bob -played -winrate```

Shows just the number of games played and win rate for Bob

## Common flags

These flags can be applied to any command

### verbose

Will log more detailed information, for example SQL queries used
e.g. 

```ego add -name=Alice -verbose```

### dbpath

Will use a different path to the sqlite database. For the ```create``` command, this will attempt to create the database at the specified path.

e.g.
```
ego create -dbpath='/home/username/egocustom.db'
ego add -name=John -dbpath='/home/username/egocustom.db'
```
