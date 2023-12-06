# Process Management System using the Matrix Protocol

## Overview
In this project, we implemented a distributed process management system (dPMS) using golang.
The main.go file contains the part of the code that creates the engines and triggers the process.

The engines folder contains The logic of the program. 
The engines/base.go file contains the communication logic. 
The other files contain the type and specific variations of machines.

In the data folder we stored the different files we needed for our test run (requirements, tasks and authentication information).

## Usage
```sh
go run main.go
```
