# RExpository

<img src="REx/src/assets/GreenRexPeas.png" width="440">

### Access directly with this link https://jaimepolop.github.io/RExpository/!!


Welcome RExpository, a collection of regular expressions, also known as regex, a sequence of characters that define a search pattern. They allow you to search for and manipulate text with precision and speed.

RExpository contains a wide variety of regular expressions, organized by category, to help you find the perfect expression. These regular expressions can be used in a variety of APIs or applications.

RExpository is a project created by Jaime Polop, in colaboration with [PEASS-ng](https://github.com/carlospolop/PEASS-ng), where everyone can contribute with additional regular expresions. We regularly update our repository with new regular expressions and updates to existing ones, so be sure to check back often for the latest additions!! 


This project was generated with [Angular CLI](https://github.com/angular/angular-cli) version 14.2.3.

## Quick Start
### regexFinder Usage

regexFinder gives the matches with a directory(or githubrepository) and the regexes, and saves the matches in a json inside ```matches``` folder.

Usage: regexFinder.go [OPTIONS]

By default regexFinder does not use the regular expresions that gives to much false positives. If you want to activate this function just type ```-c``` at the end of the execution comand.
```
Options:
  -d <local-directory>                          Find matches with a local directory you want to analyze and the regexes.
  -r <github-repository>                        Find matches with a github repository you want to analyze and the regexes.
  -rs <github-repository> <github-repository>   Find matches with several github repositories you want to analyze and the regexes.
  [-c]                                          **Always at the end!!**. Add the regex that gives false positives.
  -h                                            Show the usage of the script. 
```

### regexFinder Examples of Use

```
  go run regexFinder.go -d <local-directory> -c
  go run regexFinder.go -r https://github.com/JaimePolop/RExpository
  go run regexFinder.go -rs https://github.com/JaimePolop/RExpository https://github.com/carlospolop/PEASS-ng -c
```


## Let's improve RExpository together

Thanks you for considering contributing to the project!!

Contribute reading the **[CONTRIBUTING.md](https://github.com/JaimePolop/RExpository/blob/main/CONTRIBUTING.md)** file.
