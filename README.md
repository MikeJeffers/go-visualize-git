[![Go](https://github.com/MikeJeffers/go-visualize-git/actions/workflows/go.yml/badge.svg)](https://github.com/MikeJeffers/go-visualize-git/actions/workflows/go.yml)
![coverage](https://raw.githubusercontent.com/MikeJeffers/go-visualize-git/badges/.badges/master/coverage.svg)
# Visualize git repo statistics

Have you ever wanted to produce mediocre bar charts and graph visualizations of your git repository?  Are you tired of relying on "free" remote git hosts to produce these for you?  Have you not tried searching for any other project that does exactly what this does but a lot better?  
**Well you've come to the right place!**  

![image](https://github.com/MikeJeffers/go-visualize-git/assets/2634337/0eaa7aa6-3750-4a38-9d1d-03b4ffa67f51)

This script will parse your repositories git log and produce an underwhelming graphic with some lackluster cli args to customize the output.  
Over time there may be more features - we'll see.

Run
```sh
go run . ~/path/to/git/repo numDays mode[0|1]
```
Build
```sh
go build
```
Test
```sh
go test -v
```
