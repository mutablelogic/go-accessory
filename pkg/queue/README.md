
# Task Queue

This is an implementation of a task queue, which is a queue of tasks which are executed in order, with retry and expiry options. In order to create a task queue, use the `New` function with a [Connection Pool](../pool) and any additional options:

```go

import (
    queue "github.com/mutablelogic/go-accessory/pkg/queue"
    . "github.com/mutablelogic/go-accessory"
)

func main() {
    var pool Pool // This is the connection pool
    var opts []queue.Option // Queue options

    queue := queue.New(pool, opts...)
    if queue == nil {
        panic("Unable to create queue")
    }

    // Use the task queue
}
```

The set of options you can pass are as follows:

| Option | Description |
|--------|-------------|
| `queue.OptNamespace(string)` | Set the namespace to use for the queue |
| `queue.OptMaxAge(time.Duration)` | The maximum age for any task before expiry. By default, a task is retried without a deadline |
| `queue.OptMaxRetries(uint)` | The maximum number of retries for any task. By default, a task is retried without a maximum retry count |
| `queue.OptWorkers(uint, time.Duration)` | The number of simultaneous task workers, and the maximum time a worker is allowed to run for, or zero for no deadline. By default, the number of workers equals the number of CPU cores and there is no task deadline |
| `queue.OptBackoff(time.Duration)` | The backoff time on task failure. On first retry is made after the backoff period, the second after two times the backoff period and so forth |

Where a task queue is used on a single host, the tasks will be spread across the available workers. Where there are task queues executing on multiple hosts (using a MongoDB database, rather than sqlite), the tasks will be spread across the available workers on all hosts.

## Tasks

A task consists of:

  * A unique identifier ("Key") which is used to identify the task;
  * Optionally, a task priority. Tasks are executed in priority order, with
    higher numbers being higher priority;
  * When the task was created, and the age of the task;
  * When the task expires if it is not completed, and the number of retries
    that have be attempted to complete the task;
  * The last error that was returned when attempting to complete the task;
  * A set of "tags" to identify the parameters used to execute the task.

In order to create a new task, use the `queue.New` function:

```go
    var tags []Tag // The set of tags for the task
    task, err := queue.New(context.TODO(), tags...)
    if err != nil {
        panic("Unable to insert task")
    }
```

## Workers

You define a task worker as a function passed to the `queue.Run` method:

```go
    var worker queue.WorkerFn // The worker function

    if err := queue.Run(context.TODO(), worker); err != nil {
        panic(err)
    }
```

The signature of a worker function is `func(context.Context, Task) error`. The context should allow a worker to be cancelled and return an error for long-running tasks exceeding a deadline. The worker should return an error if a task could not be completed successfully, or `nil` if the task was completed successfully.

Your worker function can obtain task parameters using the `task.Tags()` method. The tag values are currently always strings, which can then be parsed into other types.

