<!-- ![](logo.png) -->

# Silk

* * *

A modern Source Code Management tool for multi-component software projects

* * *

**WARNING**: THIS PROJECT IS STILL IN THE EXPERIMENTAL PHASE.  
Silk is a source code management system for multi-software architectures. It provides a single project source to track, build, test, & deploy services architecture & component architecture products/projects.

## List of key features

* Version control for the project overall
* Version control on a component-by-component basis
* Coming Soon™: Interfaces to fake data for components

### About Components

Components in Silk can be local within the core project repo, but they can also be imported from other projects.
In addition, a project itself can also be imported as a component within another project.

## Usage

```shell
$ silk new my_project
$ silk n my_project
# creates a new silk project titled 'my_project'

$ silk clone
# Coming Soon™

$ silk status
$ silk s
# gets the current status of the current project or component

$ silk add my.file
# adds the file my.file to the commit buffer

$ silk add my/dir/
# Adds all files in the directory that have changes

$ silk add interactive
$ silk add i
# Coming Soon™

$ silk remove my.file
# removes the file my.file from the commit buffer

$ silk remove my/dir/
# removes all files in the directory that are in the commit buffer

$ silk remove interactive
$ silk remove i
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

$ silk log
# Coming Soon™

$ silk log 3
# Comine Soon™

$ silk log level
# Coming Soon™

$ silk log level debug
# Coming Soon™
```

### Download & Installation

```shell
# Coming Soon™
```

## Administrative Stuff

### Contributing

Keep it simple. Keep it minimal. Don't put every single feature just because you can.  
(Contribution Guide coming soon™)

### Authors or Acknowledgments

* [Ryan Pearson](https://twitter.com/LumberjackChef)

### License

This project is licensed under the MIT License (Subject to change)
