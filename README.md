<!-- ![](logo.png) -->

# Silk
================================================

A Modern VCS for Multi-Component Projects

* * *

**WARNING**: THIS PROJECT IS STILL IN THE EXPERIMENTAL PHASE.  
Silk is a version control system for multi-software architectures. It provides a single project source to track, build, test, & deploy services architecture & component architecture products/projects.

## List of key features

* Version control for the project overall
* Version control on a component-by-component basis
* Coming Soon™: Interfaces to fake data for components

## Usage

```shell
$ silk new my_project
$ silk n my_project
# creates a new silk project titled 'my_project'

$ silk status
$ silk s
# gets the current status of the current project or component

$ silk clone
# Coming Soon™

$ silk component
$ silk c
# returns a list of all components in the current project

$ silk component new my_component
$ silk c new my_component
# creates a new silk component called 'my_component'

$ silk component remove my_component
$ silk c remove my_component
# Removes a silk component called 'my_component' if it exists

$ silk version
$ silk v
# gets the version of the current project

$ silk version "0.2.0"
$ silk v "0.2.0"
# changes the version of the current project to 0.2.0
```

### Download & Installation

```shell
# Coming Soon™
```

### Contributing

Keep it simple. Keep it minimal. Don't put every single feature just because you can.

### Authors or Acknowledgments

* Ryan Pearson

### License

This project is licensed under the MIT License (Subject to change)
