# Glicko-2 Matcher
> This project implements a matcher for glicko-2 algorithm.


## Run Example

1. Download the code and enter the root directory of the `glicko2`.
2. Run matcher
   ```shell
    go test ./example -v -run Test_Matcher
   ```
3. Run settler
   ```shell
    go test ./example -v -run Test_Settler
   ```


## Hou To Use
```shell
go get -u git.17zjh.com/snake/go-pkg@feat-glicko2
```
1. Implement Player, Group, Team and Room interfaces according to your business needs.
2. Create a Macther by `NewMatcher()`, and run `matcher.Start()` to start matching.
3. When the Group starts to match, call `matcher.AddGroups(groups...)` to add the group to the matching queue and wait for the matching result.
4. When the game is over, update the `Rank` of the Team and each Player based on the result, then call `Settler.UpdateMMR(room)`.