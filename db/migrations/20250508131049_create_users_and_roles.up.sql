CREATE TABLE IF NOT EXISTS role_table (
                                          role_id UUID PRIMARY KEY,
                                          role_name VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS users_table (
                                           user_id UUID PRIMARY KEY,
                                           first_name VARCHAR(100) NOT NULL,
                                           last_name VARCHAR(100) NOT NULL,
                                           age INT NOT NULL,
                                           email VARCHAR(150) NOT NULL UNIQUE,
                                           role_id UUID NOT NULL,
                                           FOREIGN KEY (role_id) REFERENCES role_table(role_id)
);


