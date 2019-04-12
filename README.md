# GoChannels
The goal of this assignment is to parallelize reading data out of multiple files into a shared data structure using only channels for communication, while at the same time running queries against the data that has been read out so far. Queries can be of the basic GetCount to retrieve the number of times a word has been read across all files so far, or using a more complex Reduce call.

(Execution times with -askdelay=10 -infiles=data/pg1041.txt,data/pg1103.txt, using max_word reducer)
1 reader, 1 asker: ~22 Seconds
16 readers, 2 askers: ~40 Seconds
4 readers, 8 askers: ~40 Seconds
16 readers, 32 askers: Over 10 minutes
64 readers, 64 askers: Over 10 minutes

Lots of readers and writers causes contention, the program slows down incredibly fast