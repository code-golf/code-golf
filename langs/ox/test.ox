#include <oxstd.oxh>

main() {
    // Printing
    println("Hello, World!");

    // Looping
    for (decl i = 0; i < 10; ++i)
        println(i);

    // Accessing arguments
    decl args = arglist();

    for (decl i = 1; i < sizeof(args); ++i)
        println(args[i]);
}
