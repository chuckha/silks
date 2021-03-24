# Silks

Another way to manage migrations if you have defined some structs that you want stored in a database.

## Model File

Create a model file that is pure Go structs. The only imported package that is supported is `time`.

### Annotations

A model file can have annotations to change change the table name of particular models. For example, if you have a
`User` struct but want the table to be named `users`, then you would add a comment anywhere in the model file like
this: `User.tablename=users`.

Silks will reformat that file for you (TODO: right now it only prints the updated and formated version to stdout).

## Actions

`create` -> Generates create table statements
`add` -> Adds a new field to the specified go struct and a new column to the associated table.
