# goPhilosophers
A project implementing the classic Dining Philosophers Problem using Goroutines in Go.

# The Problem

The Dining Philosophers Problem is the following: there is a table with several philosophers, and each philosopher has one fork. To eat, a philosopher needs two forks it will eat in a specified amount of time. If a philosopher does not eat within a specified amount of time, they die of starvation. After eating, they sleep for a specified amount of time and then think for an unspecified amount of time. When they wake up, they must eat again before the starvation timer expires; otherwise, they die of starvation.

# Limitation

I will not use channels to make philosophers pass information to each other. I will only use them as a means of synchronization and for stopping the simulation.