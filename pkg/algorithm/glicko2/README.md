# Glicko-2 Matcher
> This project implements a matcher for glicko-2 algorithm.


## Run Example

1. Download the code and enter the root directory of the `glicko2`.
2. Run matcher
   ```shell
    go test ./example -v -run Test_Matcher
   ```


## Hou To Use
```shell
go get -u github.com/hedon954/go-matcher/pkg/algorithm/glicko2
```
1. Implement Player, Group, Team and Room interfaces according to your business needs.
2. Create a Matcher by `NewMatcher()`, and run `matcher.Start()` to start matching.
3. When the Group starts to match, call `matcher.AddGroups(groups...)` to add the group to the matching queue and wait for the matching result.
