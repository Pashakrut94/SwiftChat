CREATE TABLE IF NOT EXISTS users (    
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name VARCHAR(50),
    phone VARCHAR(12) UNIQUE,     
    password VARCHAR(150)   
)
