# No regular expressions

The entire point of creating a recursive descent parser is often to escape the performance issues and limitations of regular expressions. For example, regular expressions cannot be inlined by compiler optimization. Therefore, this package contains none of them.
