The `var` keyword has exhibits global scope and is not restricted within a block scope. So when we try to access `name` outside the `if` block it is still accessible.

`===` in js is a strict check where it checks both the value and type of the operands,whereas `==` is loose check where it only checks for value of the operands. So `5 == "5"` is true but `5 === "5"` is false

the array given has 3 elements and since js is 0 indexed the index of last element is length of array - 1 which is 2 in this case and we are trying to access the next element so it logically while js does not error out there by returning undefined stating no element is present.

`setTimeout` is a web api, asynchronous in nature which executes the logic defined in the callback after a specified time. While it does not block the main thread. Until the timer is out it stored the references of all the variables in the heap. And the setTimeout waits in the Task Queue / Callback Queue. The event loop continuously monitors the call stack and the Task Queue / Callback Queue. When the event loop finds that the call stack is empty it reaches out to the queues (Micro Task Queue has higher priority than Task Queue) and sees the setTimeout and puts it onto the callstack which is then executed and the variable is then assigned a value. In the given scenario by the time the variable is assigned a value the remaining code is already executed. The solution is to use a Promise and get the resolved data value in the .then() block. 

We are not returning any value from the function so the result is undefined, but ideally it should return the sum.

In js objects are always passed by reference so when any member is updated then it will also update the original object.

In the event listener, if you called the function it is replaced with its return value which in these case is undefined. When passing a handler to the function we shouldn't call it.

.then() is used in conjunction with a promise which captures the resolved value.
fetch("https://jsonplaceholder.typicode.com/todos") => returns a promise

    .then((res) => {                                => has a promises resolved value
        res.json();                                 => but you are returning nothing so further .then() chain has undefined
    })

    .then((data) => {                               => resolved value catched is undefined
        console.log(data);
    });

res.json() itself returns a promise so return res.json instead to capture its value in .then()

The array map method is used to map over each element of an array and to execute operations on that element and update its value by returning it. 
Ex: update the array element by 2 so we can just map over the array and add 2 to each and return it.
In given scenario, if the element is even you are returning its twice but if its odd then it's returning undefined so the undefined values in the result. Rather simply return the element to maintain it in the array.

the behavior of this changes with normal functions and arrow functions. So this is my anology, for normal functions I check "muje(fn) kisne call kiya" so if it was a normal function in this scenario then the object person is calling it so this will point to the object person.
but for arrow function I check "mere papa(any) ko kon call kar raha hai", in this scenario the papa is person object and the one calling it would be the global object, i.e. this for arrow function is pointing to the global object in this scenario and the global object does not have any field set with `name` so the value is undefined

