# RExpository

<img src="REx/src/assets/GreenRexPeas.png" width="480">

### Access directly with this link https://jaimepolop.github.io/RExpository/!!


Welcome REx, a collection of regular expressions, also known as regex, a sequence of characters that define a search pattern. They allow you to search for and manipulate text with precision and speed.

REx contains a wide variety of regular expressions, organized by category, to help you find the perfect expression. These regular expressions can be used in a variety of APIs or applications.

REx is a project created by Jaime Polop, in colaboration with [PEASS-ng](https://github.com/carlospolop/PEASS-ng), where everyone can contribute with additional regular expresions. We regularly update our repository with new regular expressions and updates to existing ones, so be sure to check back often for the latest additions!! 


This project was generated with [Angular CLI](https://github.com/angular/angular-cli) version 14.2.3.

## Quick Start
### regexFinder Usage

Usage: regexFinder.go [OPTIONS]
By default regexFinder does not use the regular expresions that gives to much false positives. If you want to activate this function just type ```-c``` at the end of the execution comand.
```
Options:
  -d <local-directory>                          Find any coincidence with a local directory you want to analyze and the regular expresions.
  -r <github-repository>                        Find any coincidence with a github repository you want to analyze and the regular expresions.
  -rs <github-repository> <github-repository>   Find any coincidence with several github repositories you want to analyze and the regular expresions.
  [-c]                                          Always at the end. Add the regular expresions that gives false positives.
  -h                                            Show the usage of the script. 
```

