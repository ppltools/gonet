### msg with colorful format
Cmsg is originally from [Masterminds/glide](https://github.com/Masterminds/glide).
Please visit the official website to get more informations.

### usage
1. log to stderr
- cmsg.Debug(format, args...)
- cmsg.Info(format, args...)
- cmsg.Warn(format, args...)
- cmsg.Err(format, args...)
- cmsg.Die(format, args...)

2. log to file
- cmsg.Stderr = filestream
- cmsg.Info(format, args...)
- ......

### License
This package is made available under an MIT-style license. See LICENSE.
