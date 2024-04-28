# ego
A bare-bones command line tool for elo tracking written in go.
Uses a local sqlite database under the hood. Only keeps track of player names and their elo, no other stats. Can only record decisive results (e.g. not draws)

## Commands
```ego create```

Initialises a local sqlite db in the location  ```~/.ego/ego.db```

```ego add -name <player name>```

Adds a player to the leaderboard with the specified name and a default starting elo of 1000. Name must be unique.

```ego show```

Displays the leaderboard

```ego record -w <winner name> -l <loser name>```

Records a win for player with name ```<winner name>``` over player with name ```<loser name>```

