# loadbalancer-algorithms

Load balancer Algorithms

## Round Robin

Round robin is a scheduling algorithm used in computing systems to manage the allocation of resources among processes or tasks. It is a pre-emptive algorithm that assigns a fixed time slice, known as a time quantum, to each process in a cyclic manner.

The basic idea of the round robin algorithm is to ensure that no process or task monopolizes the CPU, and that each process gets an equal share of CPU time. In this algorithm, each process is assigned a time quantum, which is usually a small fraction of a second. The processes are executed in a circular order, and when a process has exhausted its time quantum, it is suspended and moved to the end of the queue. The next process in the queue is then executed.

The round robin algorithm is simple and easy to implement, and it ensures fairness in the allocation of resources. However, it may not be the most efficient scheduling algorithm for all situations, as some processes may require more time than others to complete their tasks. In such cases, a priority-based scheduling algorithm may be more appropriate.

### Example

```go
package main

import (
  "fmt"
  "net/url"

  "github.com/appleboy/loadbalancer-algorithms/roundrobin"
)

func main() {
  servers := []*url.URL{
    {Host: "192.168.1.10"},
    {Host: "192.168.1.11"},
    {Host: "192.168.1.12"},
    {Host: "192.168.1.13"},
  }

  rb, err := roundrobin.New(servers...)
  if err != nil {
    panic(err)
  }

  fmt.Println(rb.NextServer().Host)
  fmt.Println(rb.NextServer().Host)
  fmt.Println(rb.NextServer().Host)
  fmt.Println(rb.NextServer().Host)
  fmt.Println(rb.NextServer().Host)
}
```

output as following

```sh
192.168.1.10
192.168.1.11
192.168.1.12
192.168.1.13
192.168.1.10
```
