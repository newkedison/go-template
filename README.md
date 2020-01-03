# go-template
This is a simple go project template manager by go mod, including:

* support windows service(via [kardianos](https://github.com/kardianos)/[service](https://github.com/kardianos/service) )
* support logging (via [logrus](https://github.com/sirupsen/logrus) )
* support a simple addin system

use [replace.sh](https://github.com/newkedison/go-template/blob/master/replace.sh) to generate deserve project type

e.g.

 <pre>
# create template for an empty project
./replace.sh new-project

# create template for an http project(via <a href="https://github.com/gin-gonic/gin">gin</a>)
./replace.sh new-project gin

# create template for an tcp server project
./replace.sh new-project tcp
</pre>
