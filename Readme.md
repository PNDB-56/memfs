memfs should do:
(1) do all CRUD for dirs and file
(2) write data to files
(3) read data from files
(4) support concurrency 

Usage:

import ("memfs")

var m = memfs.CreateInstance()
m.Pwd()
m.Ls(path) // if path is empty or null then pwd or else path
m.Mkdir(path) // ( if "/" is at beginning then its root else its curr dir )
m.rndir(name, path)

