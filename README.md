# Ticket-Server-Simuation-Go
A simulation of the Ticket Server approach to unique ID generation in Go


### Functionality
This program simulates a ticket server receiving requests for unique IDs for other processes. A ticket server is an approach to generating unique IDs in a distributed environment. The program simulated 4 concurrent processes requesting IDs from the ticket server, with 2 processes requesting IDs for a theoretical "Notes" database table, and the other 2 for a theoretical "Tasks" database table. The program runs indefinitely and prints output showing that even though the processes are concurrent, a duplicate ID will never be generated for a given table. This is just a simulation, so it is not an actual distributed environment, however the processes are meant to represent distributed systems requesting the ticket server for a unique ID.

### Background
#### Unique ID Generation in Distributed Environments
In non-distributed environments, unique ID generation can be as simple as relying on the auto-incrememnt functionality of IDs in SQL databases. However, distributed environments make unique ID generation harder because you are not always working with a single database. In the event you have many distributed databases, each with the same schema but a different collection of data, the auto-increment ID functionality can't reliably be used for unique ID generation. Below I will list some distributed environment-friendly unique ID generation methods, as well as pros and cons of each.
##### Ticket Server
Running a single server that serves other processes that request unique IDs. Due to a single server being in charge of ID generation, we can ensure there will be no concurrency issues related to unique ID generation, even in a distributed environment.
- Pros
   - Relatively easy to implement and don't need to worry about concurreny issues causing duplicate IDs.
   - Can reliably tell the order database records were created in, even across database servers.
- Cons
   - Having a single server means we have a single point of failure. If the ticket server fails, the entire system falls apart.
##### Multi-master replication
Instead of auto-incrementing by 1 each time, you increment by the number of database servers you have in your distributed system. For example, if you have 3 database servers, server 1 would generate IDs such as 1, 4, 7. Server 2 would have IDs 2, 5, 8, and server 3 would have 3, 6, 9 etc.
- Pros
   - Logically easy to understand and implement
- Cons
   - Makes it harder to scale your database server pool. If you need to increase or decrease the number database servers you have, you will need to change the auto-increment functionality to match your pool size, and you might need to remap data between servers so that the data lines up with the new auto-incrememnt pattern.
   - Cannot easily tell the order database records were created in.
##### UUID
A universally unique identifier (uuid). 128-bits long and guaranteed to be unique for practical purposes.
- Pros
   - Easy to generate with existing libraries (little engineering work needed to implement the logic for us)
   - Guaranteed to be unique for practical purposes (theoretically there could be duplicate IDs generated, but the odds are small enough to be negligible)
- Cons
   - Our design needs to accept 128-bit IDs
   - We are unable to tell the order in which the IDs were generated (we cannot easily tell if one object is older or newer than another based on ID)
##### Twitter's Snowflake Method
An algorithm that effectively takes the number of bits allocated for the IDs and splits them into subsections. The subsections can represent the ID of the server that generated the ID, the ID of the datacenter housing that server, a timestamp of time elapsed since an epoch, and a sequence number (counter for how many IDs a server has generated in a fixed amount of time. This usually resets every fixed amount of time.)
- Pros
   - A high degree of flexibility. You can set the subsection lengths to different numbers of bits based on what your needs are.
   - IDs are sortable chronologically.
- Cons
   - Can be more difficult to implement than the other algorithms listed above.

### Design Decisions

#### Using Channels to communicate between ticket server and ticket agents
I am unsure how well the current design would work at scale as I am using unbuffered channels for communication between the ticket server and ticket agents. Unbuffered channels effectively block until the communication can occur between processes, which for simulation purposes is good, but at larger scale could end up being a bottleneck if a lot of blocking is occurring. Using an unbuffered channel might help with this (having a pool of the next avaialbe unique IDs waiting for ticket agents), however this would likely require a larger refactor of the program.
