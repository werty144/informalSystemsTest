# Liars lie
## Overview
This project provides an emulation for the _Liars lie_ game which is to be played by multiple distributed agents.
## Game description
<span style="color:grey">For the detailed description please see the [discription file](https://github.com/werty144/informalSystemsTest/blob/main/task/task.pdf).</span>

In _Liars lie_ there are multiple distributed processes called _agents_ and a dedicated process called _client_. 
There is also an integer called a _network value_. The network value is known to the agents but not to the client and the
goal of the client is to deduce the network value communicating with agents.

There are agents of two types: correct and liars. Correct agents when asked respond with the actual network value 
whereas liars respond with some other integer that might differ among the liars.

## Usage 
When launched the application accepts commands via the console interface. There are 3 supported commands:
* start
  ```
  start --value v --max-value max --num-agents number --liar-ratio ratio
  ```
  Launches a number of independent agents, with number * (1-ratio) correct agents always responding
  with the specified integer value v, and (number * ratio) liar agents responding x with x != v and 1 <= x
  <= max. This command starts the game and, when ready, produces the ```agents.config``` file, which contains
  agents' TCP addresses, and prints “ready” on the terminal.
* play
  ```
  play
  ```
  Upon invocation, the client reads the ```agents.config``` file, connects to the agents, plays a
  round of the game, and prints the network value v.
* stop
  ```
  stop
  ```
  Stops all agents listed in the file ```agents.config```, removes the information about stopped agents
  from the file, and exits the executable.

### Input requirements
The program assumes that ```1 <= v, number <= 2^32 - 1```, ```2 <= max <= 2^32-1``` and that ```0 <= ratio <= 1``` and ```ratio``` is a 64-bit
precision float. In case assumptions are violated, the behavior is undefined.

## Guarantees
Note that the game has a _liars ratio_ as a parameter. It can be shown that in case _liars ratio_ >= 0.5 and _max-value_ >= 3
it is impossible to design an algorithm that would always result in the client knowing the real network value.

Hence the implementation **does not** guarantee that the client gets the right value in case _liars ratio_ >= 0.5.
Otherwise, the guarantee is given.

## Implementation details
* The application is implemented in **GO** language
* Each agent and the client run in the separate goroutine
* Client and agents communicate by the TCP
* Agents respond with their value to any received message

For more details, please take a look at the comments in the code.

## Refernces
This project is a test task for the [Informal systems](https://informal.systems/) company. 
