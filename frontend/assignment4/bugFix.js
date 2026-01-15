// function printName() {
//     let name;
//     if (true) {
//         name = "Akshay";
//     }
//     console.log(name);
// }
// printName();

// let age = 18;
// if (age === 18) {
//     console.log("Adult");
// } else {
//     console.log("Minor");
// }

// const arr = [10, 20, 30];
// for (let i = 0; i < arr.length; i++) {
//     console.log(arr[i]);
// }

// let data;
// function loadData() {
//     return new Promise((resolve) => {
//         setTimeout(() => {
//             data = "loaded"
//             resolve(data)
//         }, 1000)
//     })
// }
// loadData().then((data) => console.log(data))

// function add(a, b) {
//     if (typeof a !== "number" || typeof b !== "number") {
//         console.log("operands must be numbers")
//         return null;
//     }
//     return a + b;
// }
// const result1 = add(2, 3);
// // const result1 = add("2", "3");
// console.log(result1);

// const user = {
//     name: "John",
//     age: 25,
// };
// function updateAge(u) {
//     u.age = 30;
// }
// updateAge(user);
// console.log(user.age);

// fetch("https://jsonplaceholder.typicode.com/todos/10")
//     .then((res) => {
//         return res.json();
//     })
//     .then((data) => {
//         console.log(data);
//     })
//     .catch((err) => console.log(err.message))

// const nums = [1, 2, 3, 4];
// const result = nums.map(n => {
//     if (n % 2 === 0) {
//         return n * 2;
//     }
//     return n
// });
// console.log(result);

// const person = {
//     name: "Amar",
//     greet: function () {
//         console.log("Hello " + this.name);
//     },
// };
// person.greet();