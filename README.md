# ego
A bare-bones command line tool for elo tracking written in go.
Uses a local sqlite database under the hood.

## Commands
```ego create```

Initialises a local sqlite db in the location  ```~/.ego/ego.db```

```ego add -name <player name>```

Adds a player to the leaderboard with the specified name and a default starting elo of 1000. Name must be unique.

```ego add -name <player name> -elo <elo>```

Adds a player to the leaderboard with the specified name and starting elo

```ego leaderboard```

Displays the leaderboard

```ego games```

Displays all previous games (newest to oldest)

```ego games -limit <n>```


Displays the ```n``` most recent games

```ego games -player <player name>```

Displays all previous games involving one player

```ego record -p1 <name of player 1> -p2 <name of player 2> 11-3```

Records an 11-3 win for player 1 over player 2. Note that if the loser scores 0 points, the elo change will be doubled

