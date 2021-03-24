# Silks

Another way to manage migrations.

## When and how to use this project

`silks` works best when a project can put all the models of the system that need to be represented in SQL into one file.

This works particularly well if the project is following an architecture that uses dependency inversion. At some point,
the domain models will become translated into models that will fit into a database (and vice versa).
The file that defines those types is what this project is for. It helps keep the database in sync with the Go file.

## Model File

Create a model file that contains only pure Go structs. The only imported package that is supported is `time`.

### Annotations

A model file can have annotations to change change the table name of particular models. For example, if you have a
`User` struct but want the table to be named `users`, then you would add a comment to the User struct like
this: `User.tablename=users`.

Silks will reformat that file (TODO: right now it only prints the updated and formated version to stdout).

## Actions

`create` -> Generates create table statements

`add` -> Adds a new field to the specified go struct and a new column to the associated table.

`rename` -> Renames an existing field on a model to a new name.
