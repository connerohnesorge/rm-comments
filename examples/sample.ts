// TypeScript sample file
interface User {
  name: string;
  // User age
  age: number;
}

/* Multi-line comment
   explaining the function
*/
function greet(user: User): string {
  // Return greeting message
  return `Hello, ${user.name}!`; // inline comment
}

export { greet };
