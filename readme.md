# Neon Dream Runner

### Overview
Neon Dream Runner is the interpreter for the Neon programming language written in ***Go***. It is built to bring Neon to life. The interpreter is designed to be efficient, robust, and aligned with Neon's philosophy of flexibility and power.

### Utilities

1. Interactive command prompt for interaction with the language:  
```go run main.go```

2. Script execution:  
```go run main.go file.ne```

## About Neon
Neon is a general-purpose programming language with an adaptable level of abstraction, oriented by events and aspects, and featuring a light and clean syntax, combining the best of the imperative and functional worlds.

Our proposal is to provide access to tools not present in conventional languages, such as the choice of abstraction level, allowing low-level access, including control or shutdown of the garbage collector.

For more information, please visit our official reference repository: [Official Repository](https://github.com/ToniLommez/Neon) 

## Package Layout

1. ***pkg***  
Core library used for building the interpreter, divided into packages that abstract each phase of the process.

2. ***doc***  
Documentation of the language used for building the interpreter.

3. ***cmd***  
Other drivers used in the project complementary to the core library.

4. ***main.go***  
The main file used to run the interpreter.
