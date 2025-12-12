// Rust sample file
use std::fmt;

/// This is a doc comment
pub struct Person {
    // Person's name
    name: String,
    age: u32, // Person's age
}

/* Multi-line comment
   for the implementation
*/
impl Person {
    // Constructor
    pub fn new(name: String, age: u32) -> Self {
        Person { name, age }
    }
}

fn main() {
    // Create a person
    let person = Person::new("Alice".to_string(), 30);
    println!("Hello, {}!", person.name); // Print greeting
}
