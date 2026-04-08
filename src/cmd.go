package main

/**
create dir, file
read dir, file
update/rename dir, file
Delete dir, file
*/

type Command struct {
	cmd    string
	flags  []string
	params []string
}
