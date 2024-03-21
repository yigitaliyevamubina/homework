CREATE TABLE admins (
    id UUID PRIMARY KEY NOT NULL,
    full_name VARCHAR(200) NOT NULL,
    age INT,
    email TEXT NOT NULL,
    username VARCHAR(200) NOT NULL,
    password TEXT NOT NULL,
    role VARCHAR(100) NOT NULL,
    refresh_token TEXT
    );

INSERT INTO admins (id, full_name, age, username, email, password, role, refresh_token)
VALUES (
            'e74a31c2-ade8-444c-8aa2-4cd644d9db8f',
            'Super',
            20,
            'a',
            'mubina@gmail.com',
            '$2a$14$TiUQ5f5b9S/R7yMa5YdVoO5eCqp6sxkeo0RWNN7dC.vdHUleIahnq',
            'superadmin',
            'refresh_token'
        );
