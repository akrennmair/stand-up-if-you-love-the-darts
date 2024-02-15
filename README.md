# [stand-up-if-you-love-the-darts](https://www.youtube.com/watch?v=TfwnO3T5TY8)

Four programs for calculating darts checkouts given a total to achieve and a number of darts to achieve it, bearing in mind you have to finish on a double or the bull.

The programs are almost identical except one is written in Python (which one takes an hour to find all of the nine dart checkout), one is written in C (which takes less than a minute), one is written in Fortran (also less than a minute), and two slightly different versions of essentially the same Go program, one that is a virtually identical translation of the C version to Go, while the second one is slightly nicer without global variables (the Go versions are slightly slower than the C version).

They all have one small optimisation to reduce the number of darts you have to search for by calculating the minimum possible dart to reach a score, assuming all of the other darts scored their maximums.

A significant speed up could be had from caching, however I can't think of a way to do it while also keeping track of the darts in the checkout.

## C program

Compile using your compiler of choice. I used GCC with the following command:
`gcc -o check_me_out check_me_out.c -O3`

The program can then be run with
`./check_me_out [target_score] [number_of_darts]`

I don't do any checking on inputs so don't put any naughty values in.

When I run it with a target score of 501 and 9 darts I get the expected 3944 possible cominations in a time of just over 30 seconds.

## Fortran program

I used GFORTRAN 13.1.0 to compile with the command:
`gfortran -O3 DartCheckouts.f90 -o DartCheckouts.exe`

and ran the program with
`./DartCheckouts.exe [target score] [number_of_darts]`

With a target score of 501 and 9 darts I get the expected 3944 possible combinations in just over 26 seconds.

Also had a go at multithreading, added openMP directives to code and compiled with :
`gfortran -fopenmp -O3 DartCheckouts.f90 -o DartCheckouts.exe`

which runs for the same conditions in 7s.

## Python Program

Run using `python check_me_out.py [target_score] [number_of_darts]`, again I've done no error checking so no naughty inputs please.

It's written in pure python but does use a walrus operator so at least Python 3.8 and up, also I think some of the type hinting uses a style that was only introduced in 3.10, so let's say that to be safe.

On my machine this program runs in 2 hours 45 minutes.

## Go Program(s)

From either `go/go_like_c` or `go/go_nicer`, just run `go run check_me_out.go`.

## Graph of possible checkouts

For fun here is a graph of how many combinations there are for each possible 3 dart checkout.

![checkouts](https://github.com/fred-cook/stand-up-if-you-love-the-darts/assets/135046797/1d3ec348-187c-4848-a4b4-1379cefbf7b5)
