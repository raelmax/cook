# cook
A simple command line application to create projects from templates.

[![GoDoc](https://godoc.org/github.com/raelmax/cook?status.svg)](https://godoc.org/github.com/raelmax/cook)

## Installing


### Using

Get the binary for your platform: [linux](https://github.com/raelmax/cook/raw/master/bin/linux/cook) | [mac](https://github.com/raelmax/cook/raw/master/bin/darwin/cook)

Execute:
```
$ ./cook raelmax/cook-basic-template
```

This command will ask you some questions about this template and generate your
project based on your answers.

## Templates
Cook use the same syntax of cookiecutter to your projects templates, so, you can
use any templates from [cookiecutter list](https://github.com/audreyr/cookiecutter#available-cookiecutters).

### Create your own template
You can(and should!) create your own templates, and to help you with this task
we offer a example [at this repository](https://github.com/raelmax/cook-basic-template).