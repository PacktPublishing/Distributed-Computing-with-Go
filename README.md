# Distributed Computing with Go
This is the code repository for [Distributed Computing with Go](https://www.packtpub.com/application-development/distributed-computing-go?utm_source=github&utm_medium=repository&utm_campaign=9781787125384), published by [Packt](https://www.packtpub.com/?utm_source=github). It contains all the supporting project files necessary to work through the book from start to finish.
## About the Book
Distributed Computing with Go gives developers with a good idea how basic Go development works the tools to fulfill the true potential of Golang development in a world of concurrent web and cloud applications. Nikhil starts out by setting up a professional Go development environment. Then you’ll learn the basic concepts and practices of Golang concurrent and parallel development.

You’ll find out in the new few chapters how to balance resources and data with REST and standard web approaches while keeping concurrency in mind.  Most Go applications these days will run in a data center or on the cloud, which is a condition upon which the next chapter depends. There, you’ll expand your skills considerably by writing a distributed document indexing system during the next two chapters. This system has to balance a large corpus of documents with considerable analytical demands.

Another use case is the way in which a web application written in Go can be consciously redesigned to take distributed features into account. The chapter is rather interesting for Go developers who have to migrate existing Go applications to computationally and memory-intensive environments. The final chapter relates to the rather onerous task of testing parallel and distributed applications, something that is not usually taught in standard computer science curricula.

## Instructions and Navigation
All of the code is organized into folders. Each folder starts with a number followed by the application name. For example, Chapter02.



The code will look like the following:
```
package main
import (
 "fmt"
 "os"
)
func main() {
 fmt.Println(os.Getenv("NAME") + " is your uncle.")
}
```

The material in the book is designed to enable a hands-on approach. Throughout the book, a conscious effort has been made to provide all the relevant information to the reader beforehand so that, if the reader chooses, they can try to solve the problem on their own and then refer to the solution provided in the book. The code in the book does not have any Go dependencies beyond the standard library. This is done in order to ensure that the code examples provided in the book never change, and this also allows us to explore the standard library. 
The source code in the book should be placed at $GOPATH/src/distributedgo. The source code for examples given will be located inside the $GOPATH/src/distributed-go/chapterX folder, where X stands for the chapter number.
Download and install Go from https://golang.org/ and Docker from https://www.docker.com/community-edition website

## Related Products
* [Isomorphic Go](https://www.packtpub.com/web-development/isomorphic-go?utm_source=github&utm_medium=repository&utm_content=9781788394185)

* [Go Systems Programming](https://www.packtpub.com/networking-and-servers/go-systems-programming?utm_source=github&utm_medium=repository&utm_content=9781787125643)

* [Security with Go](https://www.packtpub.com/networking-and-servers/security-go?utm_source=github&utm_medium=repository&utm_campaign=9781788627917)
